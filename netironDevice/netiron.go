package netironDevice

import (
	"bufio"
	"fmt"
	"github.com/ipcjk/mlxsh/routerDevice"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
	"strings"
	"time"
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

	/* Fill our config with defaults for ssh and timesouts */
	router.GenerateDefaults(&Config)

	return &netironDevice{
		RTC: Config,
		Router: router.Router{
			PromptModes: make(map[string]string)}}
}
func (b *netironDevice) ConnectPrivilegedMode() (err error) {

	b.Router.SSHConnection, err = ssh.Dial("tcp", b.RTC.ConnectionAddr, b.RTC.SSHClientConfig)
	if err != nil {
		return err
	}

	b.Router.SSHSession, err = b.Router.SSHConnection.NewSession()
	if err != nil {
		return err
	}

	b.Router.SSHStdoutPipe, err = b.Router.SSHSession.StdoutPipe()
	if err != nil {
		return err
	}

	b.Router.SSHStdinPipe, err = b.Router.SSHSession.StdinPipe()
	if err != nil {
		return err
	}

	b.Router.SSHStdErrPipe, err = b.Router.SSHSession.StderrPipe()
	if err != nil {
		return err
	}

	err = b.Router.SSHSession.Shell()
	if err != nil {
		return err
	}

	prompt, err := b.readTill([]string{">", "#"})
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
	if err := b.write("enable\n"); err != nil {
		return false
	}

	_, err := b.readTill([]string{"Password:"})
	if err != nil {
		return false
	}

	if err := b.write(b.RTC.EnablePassword + "\n"); err != nil {
		return false
	}

	_, err = b.readTillEnabledPrompt()
	if err != nil {
		return false
	}

	b.Router.PromptMode = "sshEnabled"

	return true
}

func (b *netironDevice) write(command string) error {
	_, err := b.Router.SSHStdinPipe.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("Cant write to the ssh connection %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprintf(b.RTC.W, "Send command: %s", command)
	}
	time.Sleep(b.RTC.WriteTimeout)
	return nil
}

func (b *netironDevice) readTill(search []string) (string, error) {
	var lineBuf string
	shortBuf := make([]byte, 256)
	foundToken := make(chan struct{}, 0)
	defer close(foundToken)

WaitInput:
	for {
		/* Reset the timer, when we received bytes for reading */
		go func() {
			select {
			case <-(time.After(b.RTC.ReadTimeout)):
				if b.RTC.Debug {
					fmt.Fprint(b.RTC.W, "Timed out waiting for incoming buffer")
				}
				b.Router.SSHSession.Close()
				b.Router.SSHConnection.Close()
			case <-foundToken:
				return
			}
		}()
		var err error
		var n int
		if n, err = io.ReadAtLeast(b.Router.SSHStdoutPipe, shortBuf, 1); err != nil {
			/* FIXME, do something on EOF (could still contain our search buffer) */
			if err != io.EOF {
				return "", err
			} else if err == io.EOF {
				return "", err
			}
		}
		foundToken <- struct{}{}
		lineBuf += string(shortBuf[:n])
		for x := range search {
			if strings.Contains(lineBuf, search[x]) {
				break WaitInput
			}
		}
	}
	return string(lineBuf), nil
}

func (b *netironDevice) ConfigureTerminalMode() error {
	if err := b.write("conf t\n"); err != nil {
		return err
	}

	_, err := b.readTill([]string{"(config)#"})
	if err != nil {
		return fmt.Errorf("Cant find configure prompt: %s", err)
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Configuration mode on")
	}
	return nil
}

func (b *netironDevice) ExecPrivilegedMode(command string) error {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	if err := b.write(command + "\n"); err != nil {
		return err
	}
	_, err := b.readTillEnabledPrompt()
	if err != nil {
		return fmt.Errorf("Cant find  privileged mode: %s", err)
	}
	return nil
}

func (b *netironDevice) skipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("Cant switch to enabled mode to execute skip-page-display: %s", err)
	}

	if err := b.write("skip-page-display\n"); err != nil {
		return "", err
	}
	return b.readTill([]string{b.Router.SSHEnabledPrompt})
}

func (b *netironDevice) readTillEnabledPrompt() (string, error) {
	return b.readTill([]string{b.Router.SSHEnabledPrompt})
}

func (b *netironDevice) readTillConfigPrompt() (string, error) {
	return b.readTill([]string{b.Router.SSHConfigPrompt})
}

func (b *netironDevice) readTillConfigPromptSection() (string, error) {
	return b.readTill([]string{b.Router.SSHConfigPromptPre})
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
			if err := b.write("exit\n"); err != nil {
				return err
			}
		}
	case "sshConfig":
		if targetMode == "sshEnabled" {
			if err := b.write("end\n"); err != nil {
				return err
			}
		} else {
			if err := b.write("end\n"); err != nil {
				return err
			}
			if err := b.write("exit\n"); err != nil {
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

	if err := b.write("\n"); err != nil {
		return err
	}

	mode, err := b.readTill([]string{b.Router.SSHConfigPrompt, b.Router.SSHEnabledPrompt, b.Router.SSHUnprivilegedPrompt})
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

func (b *netironDevice) WriteConfiguration() (err error) {

	if err = b.SwitchMode("sshEnabled"); err != nil {
		return err
	}

	if err := b.write("write memory\n"); err != nil {
		return err
	}

	_, err = b.readTill([]string{"(config)#", "Write startup-config done."})
	if err != nil {
		return err
	}

	if b.RTC.Debug {
		fmt.Fprint(b.RTC.W, "Write startup-config done")
	}

	return
}

func (b *netironDevice) CloseConnection() {
	if b.Router.SSHSession != nil {
		b.Router.SSHSession.Close()
	}

	if b.Router.SSHConnection != nil {
		b.Router.SSHConnection.Close()
	}
}

func (b *netironDevice) PasteConfiguration(configuration io.Reader) (err error) {

	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	scanner := bufio.NewScanner(configuration)
	for scanner.Scan() {
		if err := b.write(scanner.Text() + "\n"); err != nil {
			return err
		}

		/* Wait till config prompt returns or not ? */
		if !b.RTC.SpeedMode {
			val, err := b.readTillConfigPromptSection()
			if err != nil {
				return err
			}
			if b.RTC.Debug {
				fmt.Fprintf(b.RTC.W, "Captured %s\n", val)
			}
		}
		fmt.Fprint(b.RTC.W, "+")
	}
	fmt.Fprint(b.RTC.W, "\n")

	return
}

func (b *netironDevice) RunCommands(commands io.Reader) (err error) {

	if err = b.SwitchMode("sshEnabled"); err != nil {
		return fmt.Errorf("Cant switch to privileged mode: %s", err)
	}

	scanner := bufio.NewScanner(commands)
	for scanner.Scan() {
		if err := b.write(scanner.Text() + "\n"); err != nil {
			return err
		}
		val, err := b.readTillEnabledPrompt()
		if err != nil && err != io.EOF {
			return err
		}
		fmt.Fprintf(b.RTC.W, "%s\n", val)
	}

	return err
}
