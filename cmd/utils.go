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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func correctSubdomainRegex(subdomain string) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?\z`)
	return r.MatchString(subdomain)
}

func checkPort(port string) bool {
	portNo, err := strconv.Atoi(port)
	if err != nil {
		reportError(fmt.Sprintf("Invalid port number %s, must be an integer between 0 and 65536", port), true)
	}
	return 0 < portNo && portNo < 65536
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
		reportError("Punch requires an SSH private key to connect to our servers.  By default we do not use your existing keypair.  "+
			"You can point to an existing key-pair by editing punch.toml or generate a single-purpose key using `punch generate-key`", false)
		return "", err
	}
	return string(buf), nil
}

func reportError(err string, exit bool) {
	if err == "" {
		fmt.Fprintf(os.Stderr, "Unexpected error occured\n")
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	if exit {
		os.Exit(1)
	}
}
