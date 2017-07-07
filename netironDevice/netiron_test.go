package netironDevice_test

import (
	"bytes"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/netironDevice"
	"strings"
	"testing"
	"time"
)

func TestLoadingSSHKey(t *testing.T) {
	_, err := netironDevice.LoadPrivateKey(strings.NewReader(sampleSSHKey))

	if err != nil {
		t.Error("Could not load Test key")
	}
}

func TestNetironConstructor(t *testing.T) {
	var Config = libhost.HostConfig{
		DeviceType:     "MLX",
		Hostname:       "localhost",
		Username:       "myuser",
		Password:       "mypassword",
		EnablePassword: "enablepassword",
	}

	router := netironDevice.NetironDevice(
		netironDevice.NetironConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if router == nil {
		t.Error("Cant create netiron object")
	}

	if router.SSHPort != 22 {
		t.Error("Wrong SSH-Port in default settings")
	}

	/*
		if Config.ReadTimeout != time.Second*15 {
			t.Error("Wrong readtimeout in default settings")
		}
	*/

	if Config.WriteTimeout != time.Second*0 {
		t.Error("Wrong writetimeout in default settings")
	}

	if router.Username != "myuser" || router.Password != "mypassword" {
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

	router := netironDevice.NetironDevice(
		netironDevice.NetironConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if router == nil {
		t.Error("Cant create netiron object")
	}

	/* Expect that
	sshClients fail
	*/
	if err := router.ConnectPrivilegedMode(); err == nil {
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

	router := netironDevice.NetironDevice(
		netironDevice.NetironConfig{HostConfig: Config, Debug: true, W: new(bytes.Buffer)})

	if err := router.DetectSetPrompt("SSH@frankfurt-rt1#"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

	if err := router.DetectSetPrompt("SSH@frankfurt-rt1>"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

}

var sampleSSHKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA3h5u/Jb0TKlwLAwOgaVeHevwMdCqwf2mJRvVMheNOeu2qSEk
18Rf3YS3URkUvZhdQmd/fafJYALamcxl1nO9IVEUvWXBIn3pjKR5Yf6rl4bl8V7n
MemWb+7vbCaSClHDsNXzn0PaDea+r3q0IbwovinRiCeLanfcctioBxiq6Z5ZXOTQ
CCczpbMXcN6r6R7cCip1+JUgsPVEuHvxngXXZYS3EbqWCXu8gmTPJGOw1k2jXgqW
OvyYoGZVtdhJO6ZMlY3gENw9oGxZIDqCjR3sx7Tb1w7fjg6zzu75EtLChxGTYc6L
TnvtZaVBvY7fKncy13W1iGKyoWXGmeEIYaGVtQIDAQABAoIBAQC5oyvVNYCGFdJn
Lcht+DyZu1fq+l/Mc+aI+yMKk3532xW1crrtDfWlGMdxIwofjxjaZ8+4wCNgd+Il
ShwOyHpYPwCbblClOCCaZ9f+266jnJ3PRibpozUU5df6Rp4lu8JWp+nNwRKcLa5O
0Ll9vFk83Yx+Q7aUTArVfVepXqdxSVSTEOp5tm2/DWLcJQeNNLDYlarfuUDQaYMc
AxQmd2AhLZEbCDvRO79qHUTpF3lPWbkfJ/YvxaU9hC6m7HyxLA5yY0yKRTY/wPEv
jhaLmJVprv5eZc0CCZ50UOszwRk8eGAujqPqsaNEqViOoNiQIGE/QABbCkXWe0sC
uRr8MLbhAoGBAPn74MgZwtU+PNdODERNSx/qInaKi/aZRdgaFLMQu7wg99YLcx1b
uAVQDd2oIAAQS1Fz+fiUU3qBKGWTRrlW6YBSffGXfoKBXOB7bld2hFoTtUY9RHc0
ekcKPiGT9bYIObQHxJqcivvgNkjEt+htH0Hd2OUVgDMsc6bPUxF9f29JAoGBAON2
4eBxl30xcOw6lY2EX2XvAnidqpBABz0RwxRskMR3ppkJCjK39wDVZHzrExhDKjV+
R/wzBEV7/NESfW3SpXCZ7mukM/fBUjMsmDRCg8A5m3+lWlm7HXk0qP3y4EKA8F6i
dMD7pk/S5Q+HMFc3kbtGP48eIZOXxRp+5gyD/ncNAoGAd14cwa/7ZtPnPXAZT2wR
GVY1yqDxoHkj7sLVa4PsATNE5MJm33fycSb+1/71+NHPBT/59wbsrayK26XtuYaU
zR+W4AvU7wBSlyaZU85V+KU8hCOxU7KNSOrNLD94rslStHKZILLrcsZnZWv53VRt
/oeukAUqSEVLnDWXltx0Q3ECgYEAo1d9gMVRec+FPb4cIxHJx9NIvQDLuOahzBLz
Obl0hAFAG2lIb3932ptim+nbPnMM3nkejFa+XH9a33Adrj20HBYOBjJWNzYWJzWA
3xZcsi8sIQ/Gv+UEl0Nfj21X6anZ8rtKiEKt/Wh+oRX9esQm3IrnnYiPqAM2wX4b
CSXIGAkCgYBsSXMQS2zV1rmW5EZ9tiH9ryckhr7HP8mBcpv0wfeGcvU59KpIB9ag
+5rJgWoRc9If9KKygdKRWF/gUu5+CDTEKYSW/2JjkB++lxw+dnhvAvUiinr2g+tv
PqhPHFa1PgSrw5rt8xsI0kjjcybwoxEQ6qxJUQQWOlI/4fvJsl8RaQ==
-----END RSA PRIVATE KEY-----
`
