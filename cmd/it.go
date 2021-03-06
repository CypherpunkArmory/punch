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
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/cypherpunkarmory/punch/tunnel"
	"github.com/reiver/go-telnet"
	"github.com/spf13/cobra"
)

type tunnelConf struct {
	port        string
	forwardType string
}

func (tc *tunnelConf) String() string {
	return fmt.Sprintf("%s:%s", tc.forwardType, tc.port)
}

var itCmd = &cobra.Command{
	Use:   "it <type:port>... [subdomain]",
	Short: "Expose local servers running on the ports you specify",
	Long: "Expose local servers running on the ports you specify.\n" +
		"Example: `punch it http:8080 https:8443 tcp:2000` will expose a local web server running on port 8080,\n" +
		"          an https web server running on port 8443 and a tcp server running on port 2000.\n" +
		"You can provide an optional argument to specify the name of a reserved subdomain you want to\n" +
		"associate this with.\n" +
		"Example: `punch it http:8080 https:8443 mydomain` will expose a local web server running on port 8080\n" +
		"          via \"http://mydomain.holepunch.io\" and an https web server running on port 8443 via\n" +
		"          \"https://mydomain.holepunch.io\".\n" +
		"Otherwise it will default to using a new unreserved subdomain.\n" +
		"Types supported are http, https and tcp.\n" +
		"You can have any number type:port pairs in one command.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 && correctSubdomainRegex(args[len(args)-1]) {
			subdomain = args[len(args)-1]
			args = args[:len(args)-1]
		}
		var err error
		var confs = make([]tunnelConf, len(args))
		for index, conf := range args {
			confs[index], err = getTunnelConfig(conf)
			if err != nil {
				break
			}
		}
		if err != nil {
			if args[0] == "chewie" {
				var caller telnet.Caller = telnet.StandardCaller
				telnet.DialToAndCall("towel.blinkenlights.nl:23", caller)
				os.Exit(1)
			}
			reportError("Input does not match the correct syntax type:port", true)
		}
		tunnelMultiple(confs)
	},
}

func getTunnelConfig(input string) (tunnelConf, error) {
	var output tunnelConf
	allDigits := regexp.MustCompile("[0-9]+")
	knownPorts := regexp.MustCompile("(http)|(https)|(tcp)")

	conf := strings.Split(input, ":")
	if len(conf) != 2 {
		return output, errors.New("bad input - can't determine port or protocol")
	}
	output.forwardType = conf[0]
	output.port = conf[1]

	if !knownPorts.Match([]byte(output.forwardType)) {
		return tunnelConf{}, errors.New("bad input - protocol must be http or https")
	}

	if !allDigits.Match([]byte(output.port)) {
		return tunnelConf{}, errors.New("bad input - port must be numeric")
	}

	return output, nil
}

func init() {
	rootCmd.AddCommand(itCmd)
}

func tunnelMultiple(confs []tunnelConf) {
	var tunnelConfigs = make([]tunnel.Config, len(confs))
	protocol := make([]string, len(confs))
	if subdomain != "" && !correctSubdomainRegex(subdomain) {
		reportError("Invalid Subdomain", true)
	}

	publicKey, err := getPublicKey(publicKeyPath)
	if err != nil {
		os.Exit(3)
	}

	for index, t := range confs {
		protocol[index] = t.forwardType
	}

	response, err := restAPI.CreateTunnelAPI(subdomain, publicKey, protocol)
	if err != nil {
		reportError(err.Error(), true)
	}

	if subdomain == "" {
		subdomain, _ = restAPI.GetSubdomainName(response.Subdomain.ID)
	}

	for index, conf := range confs {
		if !checkPort(conf.port) {
			reportError("Port is not in range[1-65535]", true)
			err := restAPI.DeleteTunnelAPI(subdomain)
			if err != nil {
				reportError("Could not delete tunnel. Use punch cleanup "+subdomain, true)
			}
		}

		connectionURL, err := url.Parse(sshEndpoint)
		if err != nil {
			reportError("The ssh endpoint is not a valid URL", true)
			os.Exit(3)
		}

		baseURL, err := url.Parse(baseURL)
		if err != nil {
			reportError("The base url is not a valid URL", true)
		}

		tunnelConfigs[index] = tunnel.Config{
			ConnectionEndpoint: *connectionURL,
			RestAPI:            restAPI,
			TunnelEndpoint:     response,
			EndpointType:       conf.forwardType,
			PrivateKeyPath:     privateKeyPath,
			EndpointURL:        *baseURL,
			LocalPort:          conf.port,
			Subdomain:          subdomain,
			LogLevel:           logLevel,
			TCPPorts:           response.TCPPorts,
		}
	}
	tunnel.StartReverseTunnel(tunnelConfigs...)
}
