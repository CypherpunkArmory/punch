package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func checkSubdomain(subdomain string) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?\z`)
	return r.MatchString(subdomain)
}

func checkPort(port int) bool {
	return 0 < port && port < 65536
}

func fixFilePath(path string) string {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}
	return path
}

func getPublicKey(path string) (string, error) {
	path = fixFilePath(path)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Unable to find public key. Either set correct path in .punch.toml or generate a key using `punch generate-key`")
		return "", err
	}
	return string(buf), nil
}

func printError(err error) {
	if err.Error() == "" {
		fmt.Fprintf(os.Stderr, "Unexpected error occured\n")
		return
	}
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}
