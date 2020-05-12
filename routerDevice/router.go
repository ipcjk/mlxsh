package router

import (
	"bufio"
	"fmt"
	"github.com/ipcjk/mlxsh/libssh"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ipcjk/mlxsh/libhost"
	"golang.org/x/crypto/ssh"
)

/*RunTimeConfig is an init struct that can be used
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

/*Router is a struct that will be used from the inside final router object */
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

	/* Command-Rewriter for general commands, e.g. 'sc:show_log => show logging' */
	CommandRewrite map[string]string

	/* Manage the SSH connection */
	SSHConnection *ssh.Client
	SSHSession    *ssh.Session
	SSHStdinPipe  io.WriteCloser
	SSHStdoutPipe io.Reader
	SSHStdErrPipe io.Reader
	ReadMsg       string
}

/* All kind of "generic" routines, than can be used directly or indirectly by our routers */

/*SetupSSH will create a tcp
connection and open a remote shell, optional a full
pseudo terminal */
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

/*ReadTill is our main function for reading data from the input SSH-channel and will
read from the input reader, till it finds a given string, else it will run into timeout
and close the SSH channel and session
*/
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
				ro.Close()
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

/*GenerateDefaults is setting default values for a RunTime-Configuration. It will also set all
SSH-client options like password or key authentication */
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
	/* Add default ciphers */
	sshClientConfig.SetDefaults()
	/* Add old ciphers for older Ironware switches */
	sshClientConfig.Ciphers = append(sshClientConfig.Ciphers, "aes128-cbc", "aes256-cbc", "3des-cbc")

	sshClientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	/* Use our private key if given on command-line */
	/* Allow authentication with ssh dsa or rsa key */
	if config.KeyFile != "" {
		if file, err := os.Open(config.KeyFile); err != nil {
			if config.Debug {
				fmt.Fprintf(config.W, "Cant load private key for ssh auth :(%s)\n", err)
			}
		} else {
			if privateKey, err := libssh.LoadPrivateKey(file); err != nil {
				if config.Debug {
					fmt.Fprintf(config.W, "Cant load private key for ssh auth :(%s)\n", err)
				}
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

/*ReadTillEnabledPrompt internal calls ReadTill, looking for the SSH enabled prompt string */
func (ro *Router) ReadTillEnabledPrompt(rtc RunTimeConfig) (string, error) {
	return ro.ReadTill(rtc, []string{ro.SSHEnabledPrompt})
}

/*ReadTillConfigPrompt internal calls ReadTill, looking for the SSH configuration prompt string */
func (ro *Router) ReadTillConfigPrompt(rtc RunTimeConfig) (string, error) {
	return ro.ReadTill(rtc, []string{ro.SSHConfigPrompt})
}

/*ReadTillConfigPromptSection internal calls ReadTill, looking for the SSH configuration
section prompt string */
func (ro *Router) ReadTillConfigPromptSection(rtc RunTimeConfig) (string, error) {
	return ro.ReadTill(rtc, []string{ro.SSHConfigPromptPre})
}

/* Write
takes a runtimeconfiguration and a command string and will write the command string into
the SSH input stream
*/
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

/*PasteConfiguration takes a runtimeconfiguration and a reader as argument. It will read the
reader line-by-line and inject configuration statements
*/
func (ro *Router) PasteConfiguration(rtc RunTimeConfig, configuration io.Reader) (err error) {
	scanner := bufio.NewScanner(configuration)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}

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

/*GetPromptMode will check and set the current prompt situation */
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

/*RunCommands will run commands inside the devices exec- or privileged mode.
Command will be read from a io.reader. */
func (ro *Router) RunCommands(rtc RunTimeConfig, commands io.Reader) (err error) {

	scanner := bufio.NewScanner(commands)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		/* Does the command start with mlxsh_? Then guess it is a command with replacement characters?*/
		if strings.HasPrefix(line, "mlxsh_") {
			/* stupid, but works, loop all rewrites and replace all patterns */
			for k, v := range ro.CommandRewrite {
				line = strings.ReplaceAll(line, k, v)
			}
		}

		if err := ro.Write(rtc, line+"\n"); err != nil {
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

/*DetectPrompt will try to detect the initial prompt and from this information will build a map of
future possible prompts, e.g. the configuration prompt. */
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

/*Close will close the SSH-session and the SSH-tcp-connection */
func (ro *Router) Close() {
	if ro.SSHSession != nil {
		ro.SSHSession.Close()
	}

	if ro.SSHConnection != nil {
		ro.SSHConnection.Close()
	}
}
