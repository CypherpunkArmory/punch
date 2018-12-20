package utilities

import (
	"io/ioutil"
	"regexp"
)

//CheckSubdomain checks if subdomain is valid
func CheckSubdomain(subdomain string) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?\z`)
	return r.MatchString(subdomain)
}

//CheckPort checks if port is in correct range
func CheckPort(port int) bool {
	return 0 < port && port < 65535
}

func GetPublicKey(keyPath string) (string, error) {
	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
