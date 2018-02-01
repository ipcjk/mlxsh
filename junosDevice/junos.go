package junosDevice

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/ipcjk/mlxsh/routerDevice"
)

type junosDevice struct {
	RTC router.RunTimeConfig
	router.Router
}

/*
JunosDevice returns a new
junosDevice object, has a init struct of type Router.RunTimeConfig
*/
func JunosDevice(Config router.RunTimeConfig) *junosDevice {
	var configureErrors = `(?i)(invalid command|unknown|Warning|Error|not found)`

	/* Fill our config with defaults for ssh and timesouts */
	router.GenerateDefaults(&Config)

	Config.ReadTimeout = time.Second * 15

	return &junosDevice{
		RTC: Config,
		Router: router.Router{
			PromptModes:        make(map[string]string),
			ErrorMatches:       regexp.MustCompile(configureErrors),
			PromptDetect:       `[@?\.\d\w-]+> ?$`,
			PromptReadTriggers: []string{">"},
			PromptReplacements: map[string][]string{
				"SSHConfigPrompt":    {">", "#"},
				"SSHConfigPromptPre": {">", "#"}},
		}}
}

func (b *junosDevice) Connect() (err error) {

	if err = b.Router.SetupSSH(b.RTC.ConnectionAddr, b.RTC.SSHClientConfig, false); err != nil {
		return err
	}

	/* JunOS  uses `> ` for prompt */
	prompt, err := b.Router.ReadTill(b.RTC, b.PromptReadTriggers)
	if err := b.DetectSetPrompt(prompt); err != nil {
		return fmt.Errorf("Detect prompt: %s", err)
	}

	if _, err = b.skipPageDisplayMode(); err != nil {
		return err
	}

	if err = b.GetPromptMode(b.RTC); err != nil {
		return
	}

	return
}

func (b *junosDevice) DetectSetPrompt(prompt string) error {
	return b.DetectPrompt(b.RTC, prompt)

}

func (b *junosDevice) write(command string) error {
	return b.Write(b.RTC, command)
}

func (b *junosDevice) ConfigureTerminalMode() error {
	if err := b.write("edit\n"); err != nil {
		return err
	}

	_, err := b.ReadTill(b.RTC, []string{"[edit]"})
	if err != nil {
		return fmt.Errorf("Cant find configure prompt: %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Configuration mode on")
	}
	return nil
}

func (b *junosDevice) ExecPrivilegedMode(command string) error {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	if err := b.write(command + "\n"); err != nil {
		return err
	}
	_, err := b.ReadTillEnabledPrompt(b.RTC)
	if err != nil {
		return fmt.Errorf("Cant find  privileged mode: %s", err)
	}
	return nil
}

func (b *junosDevice) skipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("Cant switch to enabled mode to execute terminal-length: %s", err)
	}

	if err := b.write("set cli screen-length 0\n"); err != nil {
		return "", err
	}

	return b.ReadTill(b.RTC, []string{b.SSHEnabledPrompt})
}

func (b *junosDevice) SwitchMode(targetMode string) error {

	if b.PromptMode == targetMode {
		return nil
	}

	switch b.PromptMode {
	case "sshEnabled":
		if targetMode == "sshConfig" {
			if err := b.ConfigureTerminalMode(); err != nil {
				return err
			}
		}
	case "sshConfig":
		if targetMode == "sshEnabled" {
			if err := b.write("exit configuration-mode\n"); err != nil {
				return err
			}
		}
	}
	b.PromptMode = targetMode

	return nil
}

func (b *junosDevice) rollback() (err error) {
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	if err := b.write("rollback\n"); err != nil {
		return err
	}

	_, err = b.ReadTill(b.RTC, []string{"load complete"})
	if err != nil {
		return fmt.Errorf("Cant rollback configuration after failed commit %s", err)
	}
	return
}

func (b *junosDevice) CommitConfiguration() (err error) {
	/* Juniper needs a commit  */
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	if err := b.Write(b.RTC, "commit comment \"mlxsh change\" and-quit\n"); err != nil {
		return err
	}

	_, err = b.ReadTill(b.RTC, []string{"Exiting configuration mode"})
	if err != nil {
		return fmt.Errorf("Commit not completed or not successful: %s", err)
	}

	/* give Juniper 1 second to settle down */
	time.Sleep(time.Millisecond * 1000)

	return
}

func (b *junosDevice) PasteConfiguration(configuration io.Reader) (err error) {

	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	return b.Router.PasteConfiguration(b.RTC, configuration)

}

func (b *junosDevice) RunCommands(commands io.Reader) (err error) {
	if err = b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	return b.Router.RunCommands(b.RTC, commands)
}

func (b *junosDevice) Close() {
	b.Router.Close()
}
