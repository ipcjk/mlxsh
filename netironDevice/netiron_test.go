package netironDevice_test

import (
	"bytes"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/netironDevice"
	"github.com/ipcjk/mlxsh/routerDevice"
	"testing"
	"time"
)

func TestNetironConstructor(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType:     "MLX",
		Hostname:       "localhost",
		Username:       "myuser",
		Password:       "mypassword",
		EnablePassword: "enablepassword",
	}

	singleRouter := netironDevice.NetironDevice(
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
		DeviceType:     "MLX",
		Hostname:       "localhost",
		Username:       "myuser",
		Password:       "mypassword",
		EnablePassword: "enablepassword",
	}

	singleRouter := netironDevice.NetironDevice(
		router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if singleRouter == nil {
		t.Error("Cant create netiron object")
	}

	/* Expect that
	sshClients fail
	*/
	if err := singleRouter.Connect(); err == nil {
		t.Error("Logged into localhost with default settings, this cant be true!")
	}

}

func TestDetectPrompt(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType:     "MLX",
		Hostname:       "localhost",
		Username:       "myuser",
		Password:       "mypassword",
		EnablePassword: "enablepassword",
	}

	singleRouter := netironDevice.NetironDevice(
		router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if err := singleRouter.DetectSetPrompt("SSH@frankfurt-rt1#"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

	if err := singleRouter.DetectSetPrompt("SSH@frankfurt-rt1>"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

}
