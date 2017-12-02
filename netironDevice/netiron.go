package netironDevice

import (
	"fmt"
	"github.com/ipcjk/mlxsh/routerDevice"
	"io"
	"regexp"
	"strings"
)

type netironDevice struct {
	RTC    router.RunTimeConfig
	Router router.Router
}

/*
NetironDevice returns a new
netironDevice object, has a init struct of type NetironConfig
*/
func NetironDevice(Config router.RunTimeConfig) *netironDevice {
	var configureErrors = `(?i)(Please first configure|invalid command|Invalid input|Warning|skipped due|Error)`

	router.GenerateDefaults(&Config)

	return &netironDevice{
		RTC: Config,
		Router: router.Router{
			PromptReadTriggers: []string{">", "#"},
			PromptModes:        make(map[string]string),
			ErrorMatches:       regexp.MustCompile(configureErrors)}}
}
func (b *netironDevice) Connect() (err error) {

	if err = b.Router.SetupSSH(b.RTC.ConnectionAddr, b.RTC.SSHClientConfig, false); err != nil {
		return err
	}

	prompt, err := b.Router.ReadTill(b.RTC, b.Router.PromptReadTriggers)
	if err != nil {
		return err
	}

	if err := b.DetectSetPrompt(prompt); err != nil {
		return err
	}

	/* Try login if promptMode is NonEnabled */
	if b.Router.PromptMode == "sshNonEnabled" && !b.loginDialog() {
		return fmt.Errorf("Cant login")
	}

	if _, err = b.skipPageDisplayMode(); err != nil {
		return err
	}

	if err = b.getPromptMode(); err != nil {
		return
	}

	return
}

func (b *netironDevice) DetectSetPrompt(prompt string) error {
	matched, err := regexp.MatchString(">$", prompt)
	if err == nil && matched {
		b.Router.PromptMode = "sshNonEnabled"
		b.Router.SSHUnprivilegedPrompt = prompt
	} else if err != nil {
		return fmt.Errorf("Cant run regexp for prompt detection, weird!")
	}

	matched, err = regexp.MatchString("#$", prompt)
	if err == nil && matched {
		b.Router.PromptMode = "sshEnabled"
		b.Router.SSHUnprivilegedPrompt = strings.Replace(prompt, "#", ">", 1)
	} else if err != nil {
		return fmt.Errorf("Cant run regexp for prompt detection, weird!")
	}

	/*
		FIXME: Need regex for replace the last one, not the first match
	*/
	b.Router.SSHEnabledPrompt = strings.Replace(b.Router.SSHUnprivilegedPrompt, ">", "#", 1)
	b.Router.SSHConfigPrompt = strings.Replace(b.Router.SSHUnprivilegedPrompt, ">", "(config)#", 1)
	b.Router.SSHConfigPromptPre = strings.Replace(b.Router.SSHUnprivilegedPrompt, ">", "(config", 1)

	b.Router.PromptModes["sshEnabled"] = b.Router.SSHEnabledPrompt
	b.Router.PromptModes["sshConfig"] = b.Router.SSHConfigPrompt
	b.Router.PromptModes["sshConfigPre"] = b.Router.SSHConfigPromptPre
	b.Router.PromptModes["sshNotEnabled"] = b.Router.SSHUnprivilegedPrompt

	if b.RTC.Debug {
		fmt.Fprintf(b.RTC.W, "Enabled:(%s)\n", b.Router.SSHEnabledPrompt)
		fmt.Fprintf(b.RTC.W, "Not-Enabled:(%s)\n", b.Router.SSHUnprivilegedPrompt)
		fmt.Fprintf(b.RTC.W, "Config:(%s)\n", b.Router.SSHConfigPrompt)
		fmt.Fprintf(b.RTC.W, "ConfigSection:(%s)\n", b.Router.SSHConfigPromptPre)
	}

	if b.Router.SSHEnabledPrompt == "" || b.Router.SSHUnprivilegedPrompt == "" {
		return fmt.Errorf("Cant detect any prompt")
	}

	return nil

}

func (b *netironDevice) loginDialog() bool {
	if err := b.Router.Write(b.RTC, "enable\n"); err != nil {
		return false
	}

	_, err := b.Router.ReadTill(b.RTC, []string{"Password:"})
	if err != nil {
		return false
	}

	if err := b.Router.Write(b.RTC, b.RTC.EnablePassword+"\n"); err != nil {
		return false
	}

	_, err = b.readTillEnabledPrompt()
	if err != nil {
		return false
	}

	b.Router.PromptMode = "sshEnabled"

	return true
}

func (b *netironDevice) ConfigureTerminalMode() error {
	if err := b.Router.Write(b.RTC, "conf t\n"); err != nil {
		return err
	}

	_, err := b.Router.ReadTill(b.RTC, []string{"(config)#"})
	if err != nil {
		return fmt.Errorf("Cant find configure prompt: %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Configuration mode on")
	}
	return nil
}

func (b *netironDevice) skipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("Cant switch to enabled mode to execute skip-page-display: %s", err)
	}

	if err := b.Router.Write(b.RTC, "skip-page-display\n"); err != nil {
		return "", err
	}
	return b.Router.ReadTill(b.RTC, []string{b.Router.SSHEnabledPrompt})
}

func (b *netironDevice) readTillEnabledPrompt() (string, error) {
	return b.Router.ReadTill(b.RTC, []string{b.Router.SSHEnabledPrompt})
}

func (b *netironDevice) readTillConfigPrompt() (string, error) {
	return b.Router.ReadTill(b.RTC, []string{b.Router.SSHConfigPrompt})
}

func (b *netironDevice) readTillConfigPromptSection() (string, error) {
	return b.Router.ReadTill(b.RTC, []string{b.Router.SSHConfigPromptPre})
}

func (b *netironDevice) SwitchMode(targetMode string) error {

	if b.Router.PromptMode == targetMode {
		return nil
	}

	switch b.Router.PromptMode {
	case "sshEnabled":
		if targetMode == "sshConfig" {
			b.ConfigureTerminalMode()
		} else {
			if err := b.Router.Write(b.RTC, "exit\n"); err != nil {
				return err
			}
		}
	case "sshConfig":
		if targetMode == "sshEnabled" {
			if err := b.Router.Write(b.RTC, "end\n"); err != nil {
				return err
			}
		} else {
			if err := b.Router.Write(b.RTC, "end\n"); err != nil {
				return err
			}
			if err := b.Router.Write(b.RTC, "exit\n"); err != nil {
				return err
			}
		}
	case "sshNotEnabled":
		if targetMode == "sshEnabled" {
			fmt.Println("LOGIN")
		} else {
			fmt.Println("LOGIN & CONF Mode")
		}
	}
	return nil
}

func (b *netironDevice) getPromptMode() error {

	if err := b.Router.Write(b.RTC, "\n"); err != nil {
		return err
	}

	mode, err := b.Router.ReadTill(b.RTC, []string{b.Router.SSHConfigPrompt, b.Router.SSHEnabledPrompt, b.Router.SSHUnprivilegedPrompt})
	if err != nil {
		return fmt.Errorf("Cant find command line mode: %s", err)
	}

	mode = strings.TrimSpace(mode)

	switch mode {
	case b.Router.PromptModes["sshEnabled"]:
		b.Router.PromptMode = "sshEnabled"
	case b.Router.PromptModes["sshConfig"]:
		b.Router.PromptMode = "sshConfig"
	case b.Router.PromptModes["sshNotEnabled"]:
		b.Router.PromptMode = "sshNotEnabled"
	default:
		b.Router.PromptMode = "unknown"
	}

	return nil
}

func (b *netironDevice) CommitConfiguration() (err error) {

	if err = b.SwitchMode("sshEnabled"); err != nil {
		return err
	}

	if err := b.Router.Write(b.RTC, "write memory\n"); err != nil {
		return err
	}

	_, err = b.Router.ReadTill(b.RTC, []string{"(config)#", "Write startup-config done."})
	if err != nil {
		return fmt.Errorf("Cant write memory flash: %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Write startup-config done")
	}

	return
}

func (b *netironDevice) PasteConfiguration(configuration io.Reader) (err error) {
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	return b.Router.PasteConfiguration(b.RTC, configuration)
}

func (b *netironDevice) RunCommands(commands io.Reader) (err error) {
	if err = b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	return b.Router.RunCommands(b.RTC, commands)
}

func (b *netironDevice) Close() {
	b.Router.Close()
}
