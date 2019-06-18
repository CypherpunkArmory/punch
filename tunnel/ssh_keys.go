// Punch CLI used for interacting with holepunch.io
// Copyright (C) 2018-2019  Orb.House, LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package tunnel

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"syscall"

	"github.com/ScaleFT/sshkeys"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func readPrivateKeyFile(path string) (ssh.AuthMethod, error) {
	log.Debugf("Parsing privatekey %s", path)

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
	// Return early if the SSH file is not password protected
	if x509.IsEncryptedPEMBlock(block) {
		return readEncryptedKey(buffer, path)
	}
	key, errParse := ssh.ParsePrivateKey(buffer)
	if errParse != nil {
		// IsEncryptedPEMBlock does not support checking OPENSSH keys
		// This is a work around that is only needed for encrypted openssh keys
		if block.Type == "OPENSSH PRIVATE KEY" {
			return readEncryptedKey(buffer, path)
		}
		return nil, errors.New("cannot parse SSH key file " + path)
	}
	return ssh.PublicKeys(key), nil
}

func readEncryptedKey(buffer []byte, path string) (ssh.AuthMethod, error) {
	fmt.Print("Your password: ")

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, errors.New("could not read your password " + err.Error())
	}
	fmt.Println()
	return readPasswordProtectedKey(buffer, bytePassword, path)
}

func readPasswordProtectedKey(buffer []byte, password []byte, path string) (ssh.AuthMethod, error) {
	key, err := sshkeys.ParseEncryptedRawPrivateKey(buffer, password)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("cannot parse SSH key file " + path)
	}
	keySigner, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return nil, errors.New("cannot get signer from SSH key file " + path)
	}
	return ssh.PublicKeys(keySigner), nil
}
