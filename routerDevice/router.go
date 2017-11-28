package router

import (
	"fmt"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/libssh"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
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
	PromptModes                                              map[string]string
	PromptMode                                               string
	SSHConfigPrompt, SSHEnabledPrompt, SSHUnprivilegedPrompt string

	SSHConfigPromptPre string
	SSHConnection      *ssh.Client
	SSHSession         *ssh.Session
	SSHStdinPipe       io.WriteCloser
	SSHStdoutPipe      io.Reader
	SSHStdErrPipe      io.Reader
	ReadMsg            string
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

	/* Check if StrictHostKeyCheck is needed */
	if config.StrictHostCheck {
		config.Hostkey = libssh.LoadHostKey("/Users/joerg/.ssh/known_hosts", config.Hostname, config.SSHIP, config.SSHPort)

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
