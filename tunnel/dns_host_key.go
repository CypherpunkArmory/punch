package tunnel

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

var errNoHostKeyFound = fmt.Errorf("sshfp: no host key found")

func dnsHostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	txtrecords, err := net.LookupTXT("api.holepunch.io")
	if err != nil {
		return err
	}
	// SHA256 checksum of key
	// TODO should also support other algos
	keyFpSHA256 := sha256.Sum256(key.Marshal())
	// TODO very naive way to validate, we should match on key type and algo
	//      and don't brute force check
	for _, entry := range txtrecords {
		sshfp := strings.Split(entry, " ")
		if len(sshfp) != 3 {
			continue
		}
		fingerPrint := sshfp[2]
		fp, err := hex.DecodeString(fingerPrint)
		if err != nil {
			continue
		}

		if bytes.Equal(fp, keyFpSHA256[:]) {
			return nil
		}
	}

	return errNoHostKeyFound
}
