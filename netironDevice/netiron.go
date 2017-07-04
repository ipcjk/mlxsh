package netironDevice

import (
	"bufio"
	"fmt"
	"github.com/ipcjk/mlxsh/libhost"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type NetironConfig struct {
	libhost.HostConfig
	Debug bool
	W     io.Writer
}

type netironDevice struct {
	NetironConfig

	promptModes                                              map[string]string
	promptMode                                               string
	sshConfigPrompt, sshEnabledPrompt, sshUnprivilegedPrompt string

	sshConfig          *ssh.ClientConfig
	sshConfigPromptPre string
	sshConnection      *ssh.Client
	sshSession         *ssh.Session
	sshStdinPipe       io.WriteCloser
	sshStdoutPipe      io.Reader
	sshStdErrPipe      io.Reader
}

/*
NetironDevice returns a new
netironDevice object
*/
func NetironDevice(Config NetironConfig) *netironDevice {
	sshConfig := &ssh.ClientConfig{User: Config.Username, Auth: []ssh.AuthMethod{ssh.Password(Config.Password)}}

	/* Allow authentication with ssh dsa or rsa key */
	if Config.KeyFile != "" {
		privateKey, err := PublicKeyFile(Config.KeyFile)
		if err != nil && Config.Debug {
			fmt.Fprintf(Config.W, "Cant load private key for ssh auth :(%s)\n", err)
		} else {
			sshConfig.Auth = append(sshConfig.Auth, privateKey)
		}
	}

	/* Set reasonable
	defaults
	*/
	if Config.SSHPort == 0 {
		Config.SSHPort = 22
	}

	if Config.ReadTimeout == 0 {
		Config.ReadTimeout = time.Second * 15
	}

	return &netironDevice{NetironConfig: Config, promptModes: make(map[string]string),
		sshConfig: sshConfig}
}

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func (b *netironDevice) ConnectPrivilegedMode() (err error) {
	b.sshConnection, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", b.Hostname, b.SSHPort), b.sshConfig)
	if err != nil {
		return err
	}

	b.sshSession, err = b.sshConnection.NewSession()
	if err != nil {
		return err
	}

	b.sshStdoutPipe, err = b.sshSession.StdoutPipe()
	if err != nil {
		return err
	}

	b.sshStdinPipe, err = b.sshSession.StdinPipe()
	if err != nil {
		return err
	}

	b.sshStdErrPipe, err = b.sshSession.StderrPipe()
	if err != nil {
		return err
	}

	err = b.sshSession.Shell()
	if err != nil {
		return err
	}

	b.sshUnprivilegedPrompt, err = b.readTill([]string{">"})
	if err != nil {
		return err
	}

	b.sshEnabledPrompt = strings.Replace(b.sshUnprivilegedPrompt, ">", "#", 1)
	b.sshConfigPrompt = strings.Replace(b.sshUnprivilegedPrompt, ">", "(config)#", 1)
	b.sshConfigPromptPre = strings.Replace(b.sshUnprivilegedPrompt, ">", "(config", 1)

	b.promptModes["sshEnabled"] = b.sshEnabledPrompt
	b.promptModes["sshConfig"] = b.sshConfigPrompt
	b.promptModes["sshConfigPre"] = b.sshConfigPromptPre
	b.promptModes["sshNotEnabled"] = b.sshUnprivilegedPrompt

	if b.Debug {
		fmt.Fprintf(b.W, "Enabled:(%s)\n", b.sshEnabledPrompt)
		fmt.Fprintf(b.W, "Not-Enabled:(%s)\n", b.sshUnprivilegedPrompt)
		fmt.Fprintf(b.W, "Config:(%s)\n", b.sshConfigPrompt)
		fmt.Fprintf(b.W, "ConfigSection:(%s)\n", b.sshConfigPromptPre)
	}

	if !b.loginDialog() {
		return fmt.Errorf("Cant login")
	}
	return
}

func (b *netironDevice) loginDialog() bool {
	if err := b.write("enable\n"); err != nil {
		return false
	}

	_, err := b.readTill([]string{"Password:"})
	if err != nil {
		return false
	}

	if err := b.write(b.EnablePassword + "\n"); err != nil {
		return false
	}

	_, err = b.readTillEnabledPrompt()
	if err != nil {
		return false
	}

	b.promptMode = "sshEnabled"

	return true
}

func (b *netironDevice) write(command string) error {
	_, err := b.sshStdinPipe.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("Cant write to the ssh connection %s", err)
	}

	if b.Debug {
		fmt.Fprintf(b.W, "Send command: %s", command)
	}
	time.Sleep(b.WriteTimeout)
	return nil
}

func (b *netironDevice) readTill(search []string) (string, error) {
	shortBuf := make([]byte, 1)
	lineBuffer := make([]byte, 0, 32)
	foundToken := make(chan struct{}, 0)
	defer close(foundToken)

WaitInput:
	for {
		/* Reset the timer, when we received at least 1 byte */
		go func() {
			select {
			case <-(time.After(b.ReadTimeout)):
				if b.Debug {
					fmt.Fprint(b.W, "Time out")
					fmt.Fprint(b.W, string(lineBuffer[:]))
				}
				b.sshSession.Close()
				b.sshConnection.Close()
			case <-foundToken:
				return
			}
		}()
		var err error
		if _, err = io.ReadAtLeast(b.sshStdoutPipe, shortBuf, 1); err != nil {
			/* FIXME, do something on EOF (could still contain our search buffer) */
			if err != io.EOF {
				return "", err
			} else if err == io.EOF {
				return "", err
			}
		}
		foundToken <- struct{}{}
		lineBuffer = append(lineBuffer, shortBuf[0])
		for x := range search {
			if strings.Contains(string(lineBuffer[:]), search[x]) {
				break WaitInput
			}
		}

	}

	return string(lineBuffer[:]), nil
}

func (b *netironDevice) ConfigureTerminalMode() error {
	if err := b.write("conf t\n"); err != nil {
		return err
	}

	_, err := b.readTill([]string{"(config)#"})
	if err != nil {
		return fmt.Errorf("Cant find configure prompt: %s", err)
	}

	if b.Debug {
		fmt.Fprint(b.W, "Configuration mode on")
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

func (b *netironDevice) SkipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("Cant switch to enabled mode to execute skip-page-display: %s", err)
	}

	if err := b.write("skip-page-display\n"); err != nil {
		return "", err
	}
	return b.readTill([]string{b.sshEnabledPrompt})
}

func (b *netironDevice) readTillEnabledPrompt() (string, error) {
	return b.readTill([]string{b.sshEnabledPrompt})
}

func (b *netironDevice) readTillConfigPrompt() (string, error) {
	return b.readTill([]string{b.sshConfigPrompt})
}

func (b *netironDevice) readTillConfigPromptSection() (string, error) {
	return b.readTill([]string{b.sshConfigPromptPre})
}

func (b *netironDevice) SwitchMode(targetMode string) error {
	if b.promptMode == targetMode {
		return nil
	}

	switch b.promptMode {
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

func (b *netironDevice) GetPromptMode() error {
	if err := b.write("\n"); err != nil {
		return err
	}

	mode, err := b.readTill([]string{b.sshConfigPrompt, b.sshEnabledPrompt, b.sshUnprivilegedPrompt})
	if err != nil {
		return fmt.Errorf("Cant find command line mode: %s", err)
	}

	mode = strings.TrimSpace(mode)

	switch mode {
	case b.promptModes["sshEnabled"]:
		b.promptMode = "sshEnabled"
	case b.promptModes["sshConfig"]:
		b.promptMode = "sshConfig"
	case b.promptModes["sshNotEnabled"]:
		b.promptMode = "sshNotEnabled"
	default:
		b.promptMode = "unknown"
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

	if b.Debug {
		fmt.Fprint(b.W, "Write startup-config done")
	}

	return
}

func (b *netironDevice) CloseConnection() {
	if b.sshConnection != nil {
		b.sshConnection.Close()
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
		if !b.SpeedMode {
			val, err := b.readTillConfigPromptSection()
			if err != nil {
				return err
			}
			if b.Debug {
				fmt.Fprintf(b.W, "Captured %s\n", val)
			}
		}
		fmt.Fprint(b.W, "+")
	}
	fmt.Fprint(b.W, "\n")

	return
}

func (b *netironDevice) RunCommandsFromReader(commands io.Reader) (err error) {

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
		fmt.Fprintf(b.W, "%s\n", val)
	}

	return err
}
