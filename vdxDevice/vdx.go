package vdxDevice

import (
	"bufio"
	"fmt"
	"github.com/ipcjk/mlxsh/libhost"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

/* VdxConfig is an init struct that can be used
to setup defaults for the VDX structure
*/
type VdxConfig struct {
	libhost.HostConfig
	Debug bool
	W     io.Writer
}

type vdxDevice struct {
	VdxConfig

	promptModes                       map[string]string
	promptMode                        string
	sshConfigPrompt, sshEnabledPrompt string

	sshClientConfig    *ssh.ClientConfig
	sshConfigPromptPre string
	sshConnection      *ssh.Client
	sshSession         *ssh.Session
	sshStdinPipe       io.WriteCloser
	sshStdoutPipe      io.Reader
	sshStdErrPipe      io.Reader
}

/*
VdxDevice returns a new
vdxDevice object, has a init struct of type VdxConfig
*/
func VdxDevice(Config VdxConfig) *vdxDevice {

	sshClientConfig := &ssh.ClientConfig{User: Config.Username, Auth: []ssh.AuthMethod{ssh.Password(Config.Password)}}
	/* Add default ciphers / hmacs */
	sshClientConfig.SetDefaults()
	/* Workaround for HostKeyCheck */
	sshClientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	/* Allow authentication with ssh dsa or rsa key */
	if Config.KeyFile != "" {
		if file, err := os.Open(Config.KeyFile); err != nil {
			if Config.Debug {
				fmt.Fprintf(Config.W, "Cant load private key for ssh auth :(%s)\n", err)
			}
		} else {
			if privateKey, err := LoadPrivateKey(file); err != nil && Config.Debug {
				fmt.Fprintf(Config.W, "Cant load private key for ssh auth :(%s)\n", err)
			} else {
				sshClientConfig.Auth = append(sshClientConfig.Auth, privateKey)
			}
		}
	}

	/* Set reasonable
	defaults
	*/
	if Config.SSHPort == 0 {
		Config.SSHPort = 22
	}

	if Config.ReadTimeout == 0 {
		Config.ReadTimeout = time.Second * 5
	}

	return &vdxDevice{VdxConfig: Config, promptModes: make(map[string]string),
		sshClientConfig: sshClientConfig}
}

/* LoadPrivateKey
loads ssh rsa or dsa private keys, is exported for testing
*/
func LoadPrivateKey(r io.Reader) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func (b *vdxDevice) ConnectPrivilegedMode() (err error) {
	var addr string

	if b.SSHIP != "" {
		addr = fmt.Sprintf("%s:%d", b.SSHIP, b.SSHPort)
	} else {
		addr = fmt.Sprintf("%s:%d", b.Hostname, b.SSHPort)
	}

	b.sshConnection, err = ssh.Dial("tcp", addr, b.sshClientConfig)
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

	/* VDX unfortunately needs a pseudo terminal */
	modes := ssh.TerminalModes{
		ssh.ECHO:  0, // Disable echoing
		ssh.IGNCR: 1, // Ignore CR on input.
	}

	/* We request a "dumb"-terminal, so we got rid of any control characters and colors */
	if err := b.sshSession.RequestPty("dumb", 80, 40, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	/* Request a shell */
	err = b.sshSession.Shell()
	if err != nil {
		return fmt.Errorf("request for shell failed: %s", err)
	}

	/* VDX always uses `# ` for prompt */
	prompt, err := b.readTill(b.sshStdoutPipe, []string{"# "})
	if err := b.DetectSetPrompt(prompt); err != nil {
		return fmt.Errorf("Detect prompt: %s", err)
	}

	return
}

func (b *vdxDevice) DetectSetPrompt(prompt string) error {

	myPrompt := regexp.MustCompile(`[\d\w-]+# ?$`).FindString(prompt)

	if myPrompt != "" {
		b.promptMode = "sshEnabled"
		b.sshEnabledPrompt = strings.TrimSpace(myPrompt)
	} else {
		return fmt.Errorf("Cant run regexp for prompt detection, weird! Found: %s", myPrompt)
	}

	b.sshConfigPrompt = strings.Replace(b.sshEnabledPrompt, "#", "(config)#", 1)
	b.sshConfigPromptPre = strings.Replace(b.sshEnabledPrompt, "#", "(config", 1)

	b.promptModes["sshEnabled"] = b.sshEnabledPrompt
	b.promptModes["sshConfig"] = b.sshConfigPrompt
	b.promptModes["sshConfigPre"] = b.sshConfigPromptPre

	if b.Debug {
		fmt.Fprintf(b.W, "Enabled:(%s)\n", b.sshEnabledPrompt)
		fmt.Fprintf(b.W, "Config:(%s)\n", b.sshConfigPrompt)
		fmt.Fprintf(b.W, "ConfigSection:(%s)\n", b.sshConfigPromptPre)
	}

	if b.sshEnabledPrompt == "" {
		return fmt.Errorf("Cant detect any prompt")
	}

	return nil

}

func (b *vdxDevice) write(command string) error {
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

func (b *vdxDevice) readTill(r io.Reader, search []string) (string, error) {
	var lineBuf string
	shortBuf := make([]byte, 256)
	foundToken := make(chan struct{}, 0)
	defer close(foundToken)

WaitInput:
	for {
		/* Reset the timer, when we received bytes for reading */
		go func() {
			select {
			case <-(time.After(b.ReadTimeout)):
				if b.Debug {
					fmt.Fprint(b.W, "Timed out waiting for incoming buffer")
					fmt.Fprintf(b.W, "Waited for %s %d", search[0], len(search[0]))
				}
				b.sshSession.Close()
				b.sshConnection.Close()
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

func (b *vdxDevice) ConfigureTerminalMode() error {
	if err := b.write("conf t\n"); err != nil {
		return err
	}

	_, err := b.readTill(b.sshStdoutPipe, []string{"(config)#"})
	if err != nil {
		return fmt.Errorf("Cant find configure prompt: %s", err)
	}

	if b.Debug {
		fmt.Fprint(b.W, "Configuration mode on")
	}
	return nil
}

func (b *vdxDevice) ExecPrivilegedMode(command string) error {
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

func (b *vdxDevice) SkipPageDisplayMode() (string, error) {
	if err := b.SwitchMode("sshEnabled"); err != nil {
		return "", fmt.Errorf("Cant switch to enabled mode to execute terminal-length: %s", err)
	}

	if err := b.write("terminal length 0\r\n"); err != nil {
		return "", err
	}
	return b.readTill(b.sshStdoutPipe, []string{b.sshEnabledPrompt})
}

func (b *vdxDevice) readTillEnabledPrompt() (string, error) {
	return b.readTill(b.sshStdoutPipe, []string{b.sshEnabledPrompt})
}

func (b *vdxDevice) readTillConfigPrompt() (string, error) {
	return b.readTill(b.sshStdoutPipe, []string{b.sshConfigPrompt})
}

func (b *vdxDevice) readTillConfigPromptSection() (string, error) {
	return b.readTill(b.sshStdoutPipe, []string{b.sshConfigPromptPre})
}

func (b *vdxDevice) SwitchMode(targetMode string) error {

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
	}
	return nil
}

func (b *vdxDevice) GetPromptMode() error {

	if err := b.write("\n"); err != nil {
		return err
	}

	mode, err := b.readTill(b.sshStdoutPipe, []string{b.sshConfigPrompt, b.sshEnabledPrompt})
	if err != nil {
		return fmt.Errorf("Cant find command line mode: %s", err)
	}

	mode = strings.TrimSpace(mode)

	switch mode {
	case b.promptModes["sshEnabled"]:
		b.promptMode = "sshEnabled"
	case b.promptModes["sshConfig"]:
		b.promptMode = "sshConfig"
	default:
		b.promptMode = "unknown"
	}

	return nil
}

func (b *vdxDevice) WriteConfiguration() (err error) {
	/* VDX does not need to write memory */
	return
}

func (b *vdxDevice) CloseConnection() {
	if b.sshSession != nil {
		b.sshSession.Close()
	}

	if b.sshConnection != nil {
		b.sshConnection.Close()
	}
}

func (b *vdxDevice) PasteConfiguration(configuration io.Reader) (err error) {

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

func (b *vdxDevice) RunCommands(commands io.Reader) (err error) {

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
