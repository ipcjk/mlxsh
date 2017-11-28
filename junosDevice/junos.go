package junosDevice

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

type junosDevice struct {
	RTC    router.RunTimeConfig
	Router router.Router
}

/*
JunosDevice returns a new
junosDevice object, has a init struct of type Router.RunTimeConfig
*/
func JunosDevice(Config router.RunTimeConfig) *junosDevice {

	/* Fill our config with defaults for ssh and timesouts */
	router.GenerateDefaults(&Config)

	Config.ReadTimeout = time.Second * 5

	return &junosDevice{
		RTC: Config,
		Router: router.Router{
			PromptModes: make(map[string]string)}}
}

func (b *junosDevice) ConnectPrivilegedMode() (err error) {

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

	/* Request a shell */
	err = b.Router.SSHSession.Shell()
	if err != nil {
		return fmt.Errorf("request for shell failed: %s", err)
	}

	/* Junus  uses `> ` for prompt */
	prompt, err := b.readTill(b.Router.SSHStdoutPipe, []string{">"})
	if err := b.DetectSetPrompt(prompt); err != nil {
		return fmt.Errorf("Detect prompt: %s", err)
	}

	if _, err = b.skipPageDisplayMode(); err != nil {
		return err
	}

	if err = b.getPromptMode(); err != nil {
		return
	}

	return
}

func (b *junosDevice) DetectSetPrompt(prompt string) error {

	myPrompt := regexp.MustCompile(`[\d\w-]+> ?$`).FindString(prompt)

	if myPrompt != "" {
		b.Router.PromptMode = "sshEnabled"
		b.Router.SSHEnabledPrompt = strings.TrimSpace(myPrompt)
	} else {
		return fmt.Errorf("Cant run regexp for prompt detection, weird! Found: %s", myPrompt)
	}

	b.Router.SSHConfigPrompt = strings.Replace(b.Router.SSHEnabledPrompt, ">", "#", 1)
	b.Router.SSHConfigPromptPre = strings.Replace(b.Router.SSHEnabledPrompt, ">", "#", 1)

	b.Router.PromptModes["sshEnabled"] = b.Router.SSHEnabledPrompt
	b.Router.PromptModes["sshConfig"] = b.Router.SSHConfigPrompt
	b.Router.PromptModes["sshConfigPre"] = b.Router.SSHConfigPromptPre

	if b.RTC.Debug {
		fmt.Fprintf(b.RTC.W, "Enabled:(%s)\n", b.Router.SSHEnabledPrompt)
		fmt.Fprintf(b.RTC.W, "Config:(%s)\n", b.Router.SSHConfigPrompt)
		fmt.Fprintf(b.RTC.W, "ConfigSection:(%s)\n", b.Router.SSHConfigPromptPre)
	}

	if b.Router.SSHEnabledPrompt == "" {
		return fmt.Errorf("Cant detect any prompt")
	}

	return nil

}

func (b *junosDevice) write(command string) error {
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

func (b *junosDevice) readTill(r io.Reader, search []string) (string, error) {
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
					fmt.Fprintf(b.RTC.W, "Waited for %s %d", search[0], len(search[0]))
				}
				b.Router.SSHSession.Close()
				b.Router.SSHConnection.Close()
			case <-foundToken:
				return
			}
		}()
		var err error
		var n int
		if n, err = io.ReadAtLeast(r, shortBuf, 1); err != nil {
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

func (b *junosDevice) ConfigureTerminalMode() error {
	if err := b.write("edit\n"); err != nil {
		return err
	}

	_, err := b.readTill(b.Router.SSHStdoutPipe, []string{"[edit]"})
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
	_, err := b.readTillEnabledPrompt()
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
	return b.readTill(b.Router.SSHStdoutPipe, []string{b.Router.SSHEnabledPrompt})
}

func (b *junosDevice) readTillEnabledPrompt() (string, error) {
	return b.readTill(b.Router.SSHStdoutPipe, []string{b.Router.SSHEnabledPrompt})
}

func (b *junosDevice) readTillConfigPrompt() (string, error) {
	return b.readTill(b.Router.SSHStdoutPipe, []string{b.Router.SSHConfigPrompt})
}

func (b *junosDevice) readTillConfigPromptSection() (string, error) {
	return b.readTill(b.Router.SSHStdoutPipe, []string{b.Router.SSHConfigPromptPre})
}

func (b *junosDevice) SwitchMode(targetMode string) error {

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
			if err := b.write("exit configuration-mode\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *junosDevice) getPromptMode() error {

	if err := b.write("\n"); err != nil {
		return err
	}

	mode, err := b.readTill(b.Router.SSHStdoutPipe, []string{b.Router.SSHConfigPrompt, b.Router.SSHEnabledPrompt})
	if err != nil {
		return fmt.Errorf("Cant find command line mode: %s", err)
	}

	mode = strings.TrimSpace(mode)

	switch mode {
	case b.Router.PromptModes["sshEnabled"]:
		b.Router.PromptMode = "sshEnabled"
	case b.Router.PromptModes["sshConfig"]:
		b.Router.PromptMode = "sshConfig"
	default:
		b.Router.PromptMode = "unknown"
	}

	return nil
}

func (b *junosDevice) rollback() (err error) {
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	if err := b.write("rollback\n"); err != nil {
		return err
	}

	_, err = b.readTill(b.Router.SSHStdoutPipe, []string{"load complete"})
	if err != nil {
		return fmt.Errorf("Cant rollback configuration after failed commit %s", err)
	}
	return
}

func (b *junosDevice) WriteConfiguration() (err error) {

	/* Juniper needs a commit  */
	if err = b.SwitchMode("sshConfig"); err != nil {
		return err
	}

	if err := b.write("commit\n"); err != nil {
		return err
	}

	_, err = b.readTill(b.Router.SSHStdoutPipe, []string{"configuration check succeeds"})
	if err != nil {
		/* Try to roll back */
		if b.rollback() != nil {
			return fmt.Errorf("Configuration check failed, rollback failed also %s", err)
		} else {
			return fmt.Errorf("Configuration check failed, rollback success %s", err)
		}
	}

	_, err = b.readTill(b.Router.SSHStdoutPipe, []string{"commit complete"})
	if err != nil {
		return fmt.Errorf("Commit not completed or not successful: %s", err)
	}

	return
}

func (b *junosDevice) CloseConnection() {
	if b.Router.SSHSession != nil {
		b.Router.SSHSession.Close()
	}

	if b.Router.SSHConnection != nil {
		b.Router.SSHConnection.Close()
	}
}

func (b *junosDevice) PasteConfiguration(configuration io.Reader) (err error) {

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

func (b *junosDevice) RunCommands(commands io.Reader) (err error) {

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
