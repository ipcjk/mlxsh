package libssh

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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

func LoadHostKey(fileName, hostname, ip string, port int) (hostKey ssh.PublicKey) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer file.Close()

	if ip != "" {
		hostname = ip
	}

	if port != 22 {
		hostname = fmt.Sprintf("[%s]:%d", hostname, port)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}

		if strings.Contains(fields[0], hostname) {
			log.Println(fields)

			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil
			}
			break
		}
	}

	if hostKey == nil {
		return nil
	}

	return hostKey
}
