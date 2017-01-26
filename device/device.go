package device

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"time"
	"io"
	"log"
	"strings"
	"bufio"
)

const (
	DEVICE_MLX = iota
)

type brocade_device struct {
	model, port int
	hostname, enable, username, password string
	readTimeout time.Duration
	writeTimeout time.Duration
	debug bool
	sshEnabledPrompt, sshConfigPrompt string
	sshSession  *ssh.Session
	sshConfig *ssh.ClientConfig
	sshStdinPipe io.WriteCloser
	sshStdoutPipe io.Reader
	sshStdErrPipe io.Reader
	sshConnection *ssh.Client
}

func Brocade (model int, hostname string, port int, enable, username, password string, readTimeout time.Duration,
writeTimeout time.Duration, debug bool)  *brocade_device  {
	return &brocade_device{model: model, port: port, hostname: hostname, enable: enable, readTimeout: readTimeout,
		writeTimeout: writeTimeout, debug: debug,
		sshConfig: &ssh.ClientConfig{User:username, Auth:[]ssh.AuthMethod{ssh.Password(password)}}}
}

func (b *brocade_device) ConnectPrivilegedMode()  {
	var err error
	b.sshConnection, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", b.hostname, b.port), b.sshConfig)
	if err != nil {
		panic(err)
	}

	b.sshSession, err = b.sshConnection.NewSession()
	if err != nil {
		panic(err)
	}

	b.sshStdoutPipe, err = b.sshSession.StdoutPipe()
	if err != nil {
		panic(err)
	}

	b.sshStdinPipe, err = b.sshSession.StdinPipe()
	if err != nil {
		panic(err)
	}

	b.sshStdErrPipe, err = b.sshSession.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = b.sshSession.Shell()
	if err != nil {
		panic(err)
	}

	capture, err := b.readTill(">")
	if err != nil {
		log.Fatal("Cant find login screen")
	}

	b.sshEnabledPrompt = strings.Replace(capture, ">", "#", 1)
	b.sshConfigPrompt = strings.Replace(capture,  ">", "(config)#", 1)

	if b.loginDialog()  && b.debug {
		log.Println("Logged in")
	}
	return
}

func (b *brocade_device) loginDialog () bool {
	b.write("enable\n")
	_, err := b.readTill("Password:")
	if err != nil {
		log.Fatal(err)
	}

	b.write(b.enable+"\n")
	_, err = b.readTillPrompt()
	if err != nil {
		log.Fatal(err)

	}

	return true
}

func (b *brocade_device) write(command string) {
	_, err := b.sshStdinPipe.Write([]byte(command))
	if err != nil {
		log.Fatal(err)

	}
	time.Sleep(b.writeTimeout)
}

func (b *brocade_device) readTill(search string)  (string, error) {
	shortBuf := make([]byte, 1)
	lineBuffer := make([]byte, 0, 32)
	foundToken := make(chan struct{}, 0)
	defer close(foundToken)

	/* Start timeout thread */
	go func() {
		select {
		case <-(time.After(b.readTimeout)):
			log.Printf("Timeout waiting for (%s)-prompt", search)
			if b.debug {
				log.Println(string(lineBuffer[:]))
			}
			b.sshSession.Close()
			b.sshConnection.Close()
			foundToken <- struct{}{}
		case <- foundToken:
			return
		}
	}()

	for {
		if _, err := io.ReadAtLeast(b.sshStdoutPipe, shortBuf, 1); err != nil {
			return string(lineBuffer[:]), err
		}
		lineBuffer = append(lineBuffer, shortBuf[0])
		if strings.Contains(string(lineBuffer[:]), search)  {
			break
		}
	}

	return string(lineBuffer[:]), nil
}

func (b *brocade_device) ConfigureTerminalMode() {
	b.write("conf t\n")
	_, err := b.readTill("(config)#")
	if err != nil {
		log.Fatal(err)
	}

	if b.debug {
		log.Println("Configuration mode on")
	}
}

func (b *brocade_device) ExecPrivilegedMode(command string) {
	b.write(command + "\n")
	_, err := b.readTillPrompt()
	if err != nil {
		log.Fatal(err)
	}
}

func (b *brocade_device) readTillPrompt() (string, error){
	return b.readTill(b.sshEnabledPrompt)
}

func (b *brocade_device) readTillConfigPrompt() (string, error){
	return b.readTill(b.sshConfigPrompt)
}

func (b *brocade_device) CloseConnection() {
	b.sshConnection.Close()
}

func (b *brocade_device) PasteConfiguration(configuration io.Reader) {
	scanner := bufio.NewScanner(configuration)
	for scanner.Scan() {
		b.write(scanner.Text()+"\n")
		val, err := b.readTillConfigPrompt()
		if err != nil {
			log.Fatal(err)
		}
		if b.debug {
			log.Printf("Last output read was: %s\n", val)
		} else {
			log.Print("+")
		}
	}
}
