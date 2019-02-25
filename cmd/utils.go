package cmd

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

//CheckSubdomain checks if subdomain is valid
func checkSubdomain(subdomain string) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?\z`)
	return r.MatchString(subdomain)
}

//CheckPort checks if port is in correct range
func checkPort(port int) bool {
	return 0 < port && port < 65536
}

//FixFilePath Fixes paths that include ~/ to full paths
func fixFilePath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	return path
}

//GetPublicKey Returns publickey as a string
func getPublicKey(path string) (string, error) {
	path = fixFilePath(path)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Unable to find public key. Either set correct path in .punch.toml or generate a key using `punch generate-key`")
		return "", err
	}
	return string(buf), nil
}
