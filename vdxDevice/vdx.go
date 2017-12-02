package vdxDevice

import (
	"fmt"
	"github.com/ipcjk/mlxsh/routerDevice"
	"io"
	"time"
)

type vdxDevice struct {
	RTC router.RunTimeConfig
	router.Router
}

/*
VdxDevice returns a new
vdxDevice object, has a init struct of type VdxConfig
*/
func VdxDevice(Config router.RunTimeConfig) *vdxDevice {

	/* Fill our config with defaults for ssh and timesouts */
	router.GenerateDefaults(&Config)

	return &vdxDevice{
		RTC: Config,
		Router: router.Router{
			PromptModes:        make(map[string]string),
			PromptDetect:       `[@?\.\d\w-]+# ?$`,
			PromptReadTriggers: []string{"# "},
			PromptReplacements: map[string][]string{
				"SSHConfigPrompt":    {"#", "(config)#"},
				"SSHConfigPromptPre": {"#", "(conf"},
			}}}
}

func (b *vdxDevice) Connect() (err error) {

	if err = b.SetupSSH(b.RTC.ConnectionAddr, b.RTC.SSHClientConfig, true); err != nil {
		return err
	}

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

func (b *vdxDevice) DetectSetPrompt(prompt string) error {
	return b.DetectPrompt(b.RTC, prompt)
}

func (b *vdxDevice) write(command string) error {
	_, err := b.SSHStdinPipe.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("Cant write to the ssh connection %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprintf(b.RTC.W, "Send command: %s", command)
	}
	time.Sleep(b.RTC.WriteTimeout)
	return nil
}

func (b *vdxDevice) ConfigureTerminalMode() error {
	if err := b.Write(b.RTC, "conf t\n"); err != nil {
		return err
	}

	_, err := b.ReadTill(b.RTC, []string{"(config)#"})
	if err != nil {
		return fmt.Errorf("Cant find configure prompt: %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Configuration mode on")
	}
	return nil
}

func (b *vdxDevice) ExecPrivilegedMode(command string) error {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	if err := b.Write(b.RTC, command+"\n"); err != nil {
		return err
	}
	_, err := b.ReadTillEnabledPrompt(b.RTC)
	if err != nil {
		return fmt.Errorf("Cant find  privileged mode: %s", err)
	}
	return nil
}

func (b *vdxDevice) skipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("Cant switch to enabled mode to execute terminal-length: %s", err)
	}

	if err := b.Write(b.RTC, "terminal length 0\r\n"); err != nil {
		return "", err
	}
	return b.ReadTill(b.RTC, []string{b.SSHEnabledPrompt})
}

func (b *vdxDevice) SwitchMode(targetMode string) error {

	if b.PromptMode == targetMode {
		return nil
	}

	switch b.PromptMode {
	case "sshEnabled":
		if targetMode == "sshConfig" {
			b.ConfigureTerminalMode()
		} else {
			if err := b.Write(b.RTC, "exit\n"); err != nil {
				return err
			}
		}
	case "sshConfig":
		if targetMode == "sshEnabled" {
			if err := b.Write(b.RTC, "end\n"); err != nil {
				return err
			}
		} else {
			if err := b.Write(b.RTC, "end\n"); err != nil {
				return err
			}
			if err := b.Write(b.RTC, "exit\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *vdxDevice) CommitConfiguration() (err error) {
	/* VDX does not need to write memory */
	return
}

func (b *vdxDevice) PasteConfiguration(configuration io.Reader) (err error) {
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	return b.Router.PasteConfiguration(b.RTC, configuration)
}

func (b *vdxDevice) RunCommands(commands io.Reader) (err error) {
	if err = b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	return b.Router.RunCommands(b.RTC, commands)
}

func (b *vdxDevice) Close() {
	b.Router.Close()
}
