package slxDevice_test

import (
	"bytes"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/routerDevice"
	"github.com/ipcjk/mlxsh/slxDevice"
	"testing"
	"time"
)

func TestSLXConstructor(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType:     "SLX",
		Hostname:       "localhost",
		Username:       "myuser",
		Password:       "mypassword",
		EnablePassword: "enablepassword",
	}

	singleRouter := slxDevice.SlxDevice(router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if singleRouter == nil {
		t.Error("Cant create SLX object")
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
		DeviceType: "SLX",
		Hostname:   "localhost",
		Username:   "user",
		Password:   "password",
		SSHPort:    9131,
	}

	singleRouter := slxDevice.SlxDevice(router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if singleRouter == nil {
		t.Error("Cant create slx object")
	}

	if err := singleRouter.Connect(); err == nil {
		t.Error("Logged into localhost with default settings, this cant be true!")
	}

}

func TestDetectPrompt(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType: "SLX",
		Hostname:   "core-10",
		Username:   "user",
		Password:   "userpassword",
		SSHPort:    9131,
	}

	singleRouter := slxDevice.SlxDevice(router.RunTimeConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if err := singleRouter.DetectSetPrompt("core-10#"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

}
