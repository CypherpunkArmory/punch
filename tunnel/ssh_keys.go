package tunnel

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func privateKeyFile(path string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("cannot read SSH key file " + path)
	}
	if len(buffer) == 0 {
		return nil, errors.New("bad key file empty file")
	}
	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, errors.New("bad key file")
	}
	if !x509.IsEncryptedPEMBlock(block) {
		key, errParse := ssh.ParsePrivateKey(buffer)
		if errParse != nil {
			return nil, errors.New("cannot parse SSH key file " + path)
		}
		return ssh.PublicKeys(key), nil
	}
	fmt.Println("Your password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, errors.New("could not read your password " + err.Error())
	}
	key, err := ssh.ParsePrivateKeyWithPassphrase(buffer, bytePassword)
	if err != nil {
		return nil, errors.New("cannot parse SSH key file " + path)
	}
	return ssh.PublicKeys(key), nil

}
