package slxDevice

import (
	"fmt"
	"io"
	"time"

	"github.com/ipcjk/mlxsh/routerDevice"
)

type slxDevice struct {
	RTC router.RunTimeConfig
	router.Router
}

/*
VdxDevice returns a new
vdxDevice object, has a init struct of type VdxConfig
*/
func SlxDevice(Config router.RunTimeConfig) *slxDevice {

	/* Fill our config with defaults for ssh and timesouts */
	router.GenerateDefaults(&Config)

	return &slxDevice{
		RTC: Config,
		Router: router.Router{
			CommandRewrite: map[string]string{
				"mlxsh_log":        "show logging",
				"mlxsh_audit":      "show logging",
				"mlxsh_chassis":    "show chassis",
				"mlxsh_route":      "show ip route",
				"mlxsh_route6":     "show ipv6 route",
				"mlxsh_include":    "include",
				"mlxsh_pipe":       "|",
				"mlxsh_route_sum":  "show ip route summary",
				"mlxsh_route6_sum": "show ipv6 route summary",
				"mlxsh_bgp":        "show ip bgp summary",
				"mlxsh_bgp6":       "show ipv6 bgp summary",
				"mlxsh_bgpn":       "show ip bgp neighbors",
				"mlxsh_bgpn6":      "show ipv6 bgp neighbors",
				"mlxsh_vlans":      "show vlan brief",
			},
			PromptModes:        make(map[string]string),
			PromptDetect:       `[@?\.\d\w-]+# ?$`,
			PromptReadTriggers: []string{"# "},
			PromptReplacements: map[string][]string{
				"SSHConfigPrompt":    {"#", "(config)#"},
				"SSHConfigPromptPre": {"#", "(conf"},
			}}}
}

func (b *slxDevice) Connect() (err error) {

	if err = b.SetupSSH(b.RTC.ConnectionAddr, b.RTC.SSHClientConfig, true); err != nil {
		return err
	}

	prompt, err := b.Router.ReadTill(b.RTC, b.PromptReadTriggers)

	if err := b.DetectSetPrompt(prompt); err != nil {
		return fmt.Errorf("detect prompt: %s", err)
	}

	if _, err = b.skipPageDisplayMode(); err != nil {
		return err
	}

	if err = b.GetPromptMode(b.RTC); err != nil {
		return
	}

	return
}

func (b *slxDevice) DetectSetPrompt(prompt string) error {
	return b.DetectPrompt(b.RTC, prompt)
}

func (b *slxDevice) write(command string) error {
	_, err := b.SSHStdinPipe.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("cant write to the ssh connection %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprintf(b.RTC.W, "Send command: %s", command)
	}
	time.Sleep(b.RTC.WriteTimeout)
	return nil
}

func (b *slxDevice) ConfigureTerminalMode() error {
	if err := b.Write(b.RTC, "conf t\n"); err != nil {
		return err
	}

	_, err := b.ReadTill(b.RTC, []string{"(config)#"})
	if err != nil {
		return fmt.Errorf("cant find configure prompt: %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Configuration mode on")
	}

	return nil
}

func (b *slxDevice) ExecPrivilegedMode(command string) error {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("cant switch to privileged mode: %s", err)
	}

	if err := b.Write(b.RTC, command+"\n"); err != nil {
		return err
	}
	_, err := b.ReadTillEnabledPrompt(b.RTC)
	if err != nil {
		return fmt.Errorf("cant find  privileged mode: %s", err)
	}
	return nil
}

func (b *slxDevice) skipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("cant switch to enabled mode to execute terminal-length: %s", err)
	}

	if err := b.Write(b.RTC, "terminal length 0\r\n"); err != nil {
		return "", err
	}
	return b.ReadTill(b.RTC, []string{b.SSHEnabledPrompt})
}

func (b *slxDevice) SwitchMode(targetMode string) error {

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
			if err := b.Write(b.RTC, "exit configuration-mode\n"); err != nil {
				return err
			}
		} else {
			if err := b.Write(b.RTC, "exit\n"); err != nil {
				return err
			}
		}
	}
	b.PromptMode = targetMode
	return nil
}

func (b *slxDevice) CommitConfiguration() (err error) {
	/* SLX needs to commit the running to startup config */

	if err = b.SwitchMode("sshEnabled"); err != nil {
		return err
	}

	if err := b.Write(b.RTC, "copy running-config startup-config\n"); err != nil {
		return err
	}

	_, err = b.ReadTill(b.RTC, []string{"continue? [y/n]"})
	if err != nil {
		return fmt.Errorf("commit not completed or not successful: %s", err)
	}

	if err := b.Write(b.RTC, "y\n"); err != nil {
		return err
	}

	_, err = b.ReadTillEnabledPrompt(b.RTC)
	if err != nil {
		return fmt.Errorf("cant find  back to privileged mode: %s", err)
	}

	/* give the SLX 1 or 2 seconds to settle down */
	time.Sleep(time.Millisecond * 2000)

	return
}

func (b *slxDevice) PasteConfiguration(configuration io.Reader) (err error) {
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	return b.Router.PasteConfiguration(b.RTC, configuration)
}

func (b *slxDevice) RunCommands(commands io.Reader) (err error) {
	if err = b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("cant switch to privileged mode: %s", err)
	}

	return b.Router.RunCommands(b.RTC, commands)
}

func (b *slxDevice) Close() {
	b.Router.Close()
}
