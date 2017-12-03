package router

import (
	"bufio"
	"fmt"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/libssh"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

/* RunConfig is an init struct that can be used
to setup defaults and values for a Router object
*/

type RunTimeConfig struct {
	libhost.HostConfig
	Debug           bool
	W               io.Writer
	ConnectionAddr  string
	SSHClientConfig *ssh.ClientConfig
	Hostkey         ssh.PublicKey
}

/* Router is a struct that will be used from the inside final router object */
type Router struct {
	/* Manage, scan and find the prompts */
	PromptModes        map[string]string
	PromptReplacements map[string][]string
	PromptMode         string

	/* Regular expression string for detecting the prompt */
	PromptDetect string
	/* Initial trigger strings for searching the prompt */
	PromptReadTriggers []string

	/* Ready configured prompt types */
	SSHConfigPrompt, SSHEnabledPrompt, SSHUnprivilegedPrompt string
	SSHConfigPromptPre                                       string

	/* ErrorMatches is a regex to scan for error messages in the configuration terminal */
	ErrorMatches *regexp.Regexp

	SSHConnection *ssh.Client
	SSHSession    *ssh.Session
	SSHStdinPipe  io.WriteCloser
	SSHStdoutPipe io.Reader
	SSHStdErrPipe io.Reader
	ReadMsg       string
}

/* All kind of "generic" routines, than can be used directly or indirectly by our routers */

/* Setup SSH connection */
func (ro *Router) SetupSSH(addr string, clientConfig *ssh.ClientConfig, requestPty bool) (err error) {
	ro.SSHConnection, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return err
	}

	ro.SSHSession, err = ro.SSHConnection.NewSession()
	if err != nil {
		return err
	}

	ro.SSHStdoutPipe, err = ro.SSHSession.StdoutPipe()
	if err != nil {
		return err
	}

	ro.SSHStdinPipe, err = ro.SSHSession.StdinPipe()
	if err != nil {
		return err
	}

	ro.SSHStdErrPipe, err = ro.SSHSession.StderrPipe()
	if err != nil {
		return err
	}

	if requestPty {
		/* Brocade VDX unfortunately needs a pseudo terminal */
		modes := ssh.TerminalModes{
			ssh.ECHO:  0, // Disable echoing
			ssh.IGNCR: 1, // Ignore CR on input.
		}

		/* We request a "dumb"-terminal, so we got rid of any control characters and colors */
		if err := ro.SSHSession.RequestPty("dumb", 80, 40, modes); err != nil {
			return fmt.Errorf("request for pseudo terminal failed: %s", err)
		}
	}

	/* Request a shell */
	err = ro.SSHSession.Shell()
	if err != nil {
		return fmt.Errorf("request for shell failed: %s", err)
	}

	return
}

func (ro *Router) ReadTill(rtc RunTimeConfig, search []string) (string, error) {
	var lineBuf string
	shortBuf := make([]byte, 512)
	foundToken := make(chan struct{}, 0)
	defer close(foundToken)

WaitInput:
	for {
		/* Reset the timer, when we received bytes for reading */
		go func() {
			select {
			case <-(time.After(rtc.ReadTimeout)):
				if rtc.Debug {
					fmt.Fprint(rtc.W, "Timed out waiting for incoming buffer")
					fmt.Fprintf(rtc.W, "Waited for %s %d", search[0], len(search[0]))
				}
				ro.SSHSession.Close()
				ro.SSHConnection.Close()
			case <-foundToken:
				return
			}
		}()
		var err error
		var n int
		if n, err = io.ReadAtLeast(ro.SSHStdoutPipe, shortBuf, 1); err != nil {
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

/* Set default values for a Configuration */
func GenerateDefaults(config *RunTimeConfig) {

	/* Set reasonable defaults */
	if config.SSHPort == 0 {
		config.SSHPort = 22
	}

	if config.ReadTimeout == 0 {
		config.ReadTimeout = time.Second * 5
	}

	/* Generate a SSH configuration profile */
	sshClientConfig := &ssh.ClientConfig{User: config.Username, Auth: []ssh.AuthMethod{ssh.Password(config.Password)}}
	/* Add default ciphers / hmacs */
	sshClientConfig.SetDefaults()
	/* Add old ciphers for older Ironware */
	sshClientConfig.Ciphers = append(sshClientConfig.Ciphers, "aes128-cbc", "3des-cbc")

	/* Check if StrictHostKeyCheck is needed */
	if config.StrictHostCheck {
		file, err := os.Open(config.KnownHosts)
		if err != nil {
			fmt.Fprintf(config.W, "Fatal: Cant open known ssh hosts file for strict ssh host authentication")
		}
		defer file.Close()

		config.Hostkey = libssh.SearchHostKey(file, config.Hostname, config.SSHIP, config.SSHPort)

		if config.Hostkey != nil {
			sshClientConfig.HostKeyCallback = ssh.FixedHostKey(config.Hostkey)
		} else {
			fmt.Fprintf(config.W, "Fatal: Cant load host key for strict ssh host authentication")
		}
	} else {
		sshClientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	/* Use our private key if given on command-line */
	/* Allow authentication with ssh dsa or rsa key */
	if config.KeyFile != "" {
		if file, err := os.Open(config.KeyFile); err != nil {
			if config.Debug {
				fmt.Fprintf(config.W, "Cant load private key for ssh auth :(%s)\n", err)
			}
		} else {
			if privateKey, err := libssh.LoadPrivateKey(file); err != nil && config.Debug {
				fmt.Fprintf(config.W, "Cant load private key for ssh auth :(%s)\n", err)
			} else {
				sshClientConfig.Auth = append(sshClientConfig.Auth, privateKey)
			}
		}
	}

	/* Copy to config object */
	config.SSHClientConfig = sshClientConfig

	/* build connectionAddr */
	if config.SSHIP != "" {
		config.ConnectionAddr = fmt.Sprintf("%s:%d", config.SSHIP, config.SSHPort)
	} else {
		config.ConnectionAddr = fmt.Sprintf("%s:%d", config.Hostname, config.SSHPort)
	}
}

func (ro *Router) ReadTillEnabledPrompt(rtc RunTimeConfig) (string, error) {
	return ro.ReadTill(rtc, []string{ro.SSHEnabledPrompt})
}

func (ro *Router) ReadTillConfigPrompt(rtc RunTimeConfig) (string, error) {
	return ro.ReadTill(rtc, []string{ro.SSHConfigPrompt})

}

func (ro *Router) ReadTillConfigPromptSection(rtc RunTimeConfig) (string, error) {
	return ro.ReadTill(rtc, []string{ro.SSHConfigPromptPre})
}

func (ro *Router) Write(rtc RunTimeConfig, command string) error {
	_, err := ro.SSHStdinPipe.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("Cant write to the ssh connection %s", err)
	}

	if rtc.Debug {
		fmt.Fprintf(rtc.W, "Send command: %s", command)
	}
	time.Sleep(rtc.WriteTimeout)
	return nil
}

func (ro *Router) PasteConfiguration(rtc RunTimeConfig, configuration io.Reader) (err error) {

	scanner := bufio.NewScanner(configuration)
	for scanner.Scan() {
		if err := ro.Write(rtc, scanner.Text()+"\n"); err != nil {
			return err
		}

		/* Wait till config prompt returns or not ? */
		if !rtc.SpeedMode {
			val, err := ro.ReadTillConfigPromptSection(rtc)
			if err != nil {
				return err
			}
			if rtc.Debug {
				fmt.Fprintf(rtc.W, "Captured %s\n", val)
			}
			if ro.ErrorMatches != nil && ro.ErrorMatches.MatchString(val) {
				return fmt.Errorf("Invalid configuration statement: %s ", scanner.Text())
			}
		}
		fmt.Fprint(rtc.W, "+")
	}
	fmt.Fprint(rtc.W, "\n")

	return
}

func (ro *Router) GetPromptMode(rtc RunTimeConfig) error {

	if err := ro.Write(rtc, "\n"); err != nil {
		return err
	}

	mode, err := ro.ReadTill(rtc, []string{ro.SSHConfigPrompt, ro.SSHEnabledPrompt})
	if err != nil {
		return fmt.Errorf("Cant find command line mode: %s", err)
	}

	mode = strings.TrimSpace(mode)

	switch mode {
	case ro.PromptModes["sshEnabled"]:
		ro.PromptMode = "sshEnabled"
	case ro.PromptModes["sshConfig"]:
		ro.PromptMode = "sshConfig"
	default:
		ro.PromptMode = "unknown"
	}

	return nil
}

func (ro *Router) RunCommands(rtc RunTimeConfig, commands io.Reader) (err error) {

	scanner := bufio.NewScanner(commands)
	for scanner.Scan() {
		if err := ro.Write(rtc, scanner.Text()+"\n"); err != nil {
			return err
		}
		val, err := ro.ReadTillEnabledPrompt(rtc)
		if err != nil && err != io.EOF {
			return err
		}
		fmt.Fprintf(rtc.W, "%s\n", val)
	}

	return err
}

func (ro *Router) DetectPrompt(rtc RunTimeConfig, prompt string) error {

	myPrompt := regexp.MustCompile(ro.PromptDetect).FindString(prompt)

	if myPrompt != "" {
		ro.PromptMode = "sshEnabled"
		ro.SSHEnabledPrompt = strings.TrimSpace(myPrompt)
	} else {
		return fmt.Errorf("Cant run regexp for prompt detection, weird! Found: %s", myPrompt)
	}

	ro.SSHConfigPrompt = strings.Replace(ro.SSHEnabledPrompt, ro.PromptReplacements["SSHConfigPrompt"][0], ro.PromptReplacements["SSHConfigPrompt"][1], 1)
	ro.SSHConfigPromptPre = strings.Replace(ro.SSHEnabledPrompt, ro.PromptReplacements["SSHConfigPromptPre"][0], ro.PromptReplacements["SSHConfigPromptPre"][1], 1)

	ro.PromptModes["sshEnabled"] = ro.SSHEnabledPrompt
	ro.PromptModes["sshConfig"] = ro.SSHConfigPrompt
	ro.PromptModes["sshConfigPre"] = ro.SSHConfigPromptPre

	if rtc.Debug {
		fmt.Fprintf(rtc.W, "Enabled:(%s)\n", ro.SSHEnabledPrompt)
		fmt.Fprintf(rtc.W, "RTC:(%s)\n", ro.SSHConfigPrompt)
		fmt.Fprintf(rtc.W, "ConfigSection:(%s)\n", ro.SSHConfigPromptPre)
	}

	if ro.SSHEnabledPrompt == "" {
		return fmt.Errorf("Cant detect any prompt")
	}

	return nil

}

func (ro *Router) Close() {
	if ro.SSHSession != nil {
		ro.SSHSession.Close()
	}

	if ro.SSHConnection != nil {
		ro.SSHConnection.Close()
	}
}
