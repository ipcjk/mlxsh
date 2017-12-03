package libssh

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
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

/* SearchHostKey searches for a host-entry in a reader  */
func SearchHostKey(r io.Reader, hostname, ip string, port int) (hostKey ssh.PublicKey) {

	var err error
	var hosts [2]string

	if ip != "" && hostname != "" {
		hosts[0] = hostname
		hosts[1] = ip
	} else if ip != "" && hostname == "" {
		hosts[0] = ip
		hosts[1] = "0xFF"
	} else {
		hosts[0] = hostname
		hosts[1] = "0xFF"
	}

	if port != 22 {
		for i := range hosts {
			hosts[i] = fmt.Sprintf("[%s]:%d", hostname, port)
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		fields := strings.Split(line, " ")
		if len(fields) != 3 {
			continue
		}

		/* found my hostname or ip? */
		if strings.Contains(fields[0], hosts[0]) || strings.Contains(fields[0], hosts[1]) {
			/* normal, old-style host format */
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err == nil {
				return hostKey
			}
		} else /* found hashed hostname? */ if strings.HasPrefix(fields[0], "|1|") {
			/* extract salt */
			_, salt, _, err := decodeHash(fields[0])
			if err == nil {
				/* try every possibility from hosts array */
				for i := range hosts {
					newHash := hashHost(hosts[i], salt)
					if encodeHash("1", salt, newHash) == fields[0] {
						hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
						if err == nil {
							return hostKey
						}
					}
				}
			}
		}

	}

	return

}

/* taken from golang crypto ssh package */
func decodeHash(encoded string) (hashType string, salt, hash []byte, err error) {

	if len(encoded) == 0 || encoded[0] != '|' {
		err = fmt.Errorf("host must start with '|'")
		return
	}
	components := strings.Split(encoded, "|")
	if len(components) != 4 {
		err = fmt.Errorf("got %d components, want 3", len(components))
		return
	}

	hashType = components[1]

	if salt, err = base64.StdEncoding.DecodeString(components[2]); err != nil {
		return "", nil, nil, fmt.Errorf("foo", err)
	}
	if hash, err = base64.StdEncoding.DecodeString(components[3]); err != nil {
		return "", nil, nil, fmt.Errorf("foo", err)
	}
	return
}

func encodeHash(typ string, salt []byte, hash []byte) string {
	return strings.Join([]string{"",
		typ,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	}, "|")
}

func hashHost(hostname string, salt []byte) []byte {
	mac := hmac.New(sha1.New, salt)
	mac.Write([]byte(hostname))
	return mac.Sum(nil)
}
