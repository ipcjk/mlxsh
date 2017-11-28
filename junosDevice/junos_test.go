package junosDevice_test

import (
	"bytes"
	"github.com/ipcjk/mlxsh/junosDevice"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/routerDevice"
	"testing"
	"time"
)

func TestJunOSConstructor(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType:     "junos",
		Hostname:       "localhost",
		Username:       "myuser",
		Password:       "mypassword",
		EnablePassword: "enablepassword",
	}

	singleRouter := junosDevice.JunosDevice(
		router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if singleRouter == nil {
		t.Error("Cant create netiron object")
	}

	if singleRouter.RTC.SSHPort != 22 {
		t.Error("Wrong SSH-Port in default settings")
	}

	if Config.WriteTimeout != time.Second*0 {
		t.Error("Wrong writetimeout in default settings")
	}

	if singleRouter.RTC.Username != "myuser" || singleRouter.RTC.Password != "mypassword" {
		t.Error("Cant match user or password")
	}

}

func TestSSHConnect(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType: "JUNOS",
		Hostname:   "localhost",
		Username:   "joerg",
		Password:   "gjh48zgu34we",
		SSHPort:    22,
	}

	singleRouter := junosDevice.JunosDevice(
		router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if singleRouter == nil {
		t.Error("Cant create junos object")
	}

	/* Expect that
	sshClients fail
	*/
	if err := singleRouter.ConnectPrivilegedMode(); err == nil {
		t.Error("Logged into localhost with default settings, this cant be true!")
	}

}

func TestDetectPrompt(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType:     "juniper",
		Hostname:       "localhost",
		Username:       "username",
		Password:       "password",
		EnablePassword: "enable",
	}

	singleRouter := junosDevice.JunosDevice(
		router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if err := singleRouter.DetectSetPrompt("joergkost@rm-core-vc-186-1>"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

}
