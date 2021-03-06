package libssh_test

import (
	"strings"
	"testing"

	"github.com/ipcjk/mlxsh/libssh"
)

func TestLoadingSSHKey(t *testing.T) {
	_, err := libssh.LoadPrivateKey(strings.NewReader(sampleSSHKey))

	if err != nil {
		t.Error("Could not load Test key")
	}
}

func TestSearchHostKey(t *testing.T) {
	var sampleHostKey = `192.168.1.76 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAvVgFjl9iLW8g8To6GCb2EhGBhtLWnPU34yJWDEziAkH6MHcOpCM/ZTHL8lOmawC1m+YfptFa5IzWtY57/EAusuUeFzT1btISVxNDSWE88HMCv4G0NT9euDtcTWThK7chxK8R1VzRMLYc3QmpxkTxFg9r+ZCMzf77SQw211y0wesu4mJ784yU025jpgmRRstUvX9AJaTo++StKbzCbakC2FjkRadWfkfnXe7KHbdi/RpY5z/4aESA41sqkHgxYS6ja1UNvP7xaSYY0lCrp5McSj3AkKCBqgqiqJPBZLgZ+B5RoBVhGS6yJ1/LRPVGJJKNcrTFuwPhC8AwacxvJpwasQ==
|1|APfqp9OCXJ3NWEk6I/TKaqpHmr4=|gJC4q5LAv84V0UeHFN0SomiRuxY= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBOY+0mKz+wOcJhE+322JTnRLBNdxWrHJMRf/S+eckJTEYTsTMEil9aBXdfiA4kSxE1rcvmWwRChcFyeNMwtSCAI=`

	hostKey := libssh.SearchHostKey(strings.NewReader(sampleHostKey), "so-call-me-maybe", "", 22)

	if hostKey == nil {
		t.Error("Hashed SSH-Key not found")
	}

	hostKey = libssh.SearchHostKey(strings.NewReader(sampleHostKey), "", "192.168.1.76", 22)

	if hostKey == nil {
		t.Errorf("Hostkey %s not found", "192.168.1.76")
	}

	hostKey = libssh.SearchHostKey(strings.NewReader(sampleHostKey), "so-call-somebody-else", "192.168.1.76", 22)

	if hostKey == nil {
		t.Errorf("Hostkey %s not found", "192.168.1.76")
	}

	hostKey = libssh.SearchHostKey(strings.NewReader(sampleHostKey), "so-call-anyone", "172.168.1.76", 2122)

	if hostKey != nil {
		t.Error("Found unknown Hostkey")
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
