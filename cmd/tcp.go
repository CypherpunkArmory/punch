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
	"net/url"
	"os"

	"github.com/cypherpunkarmory/punch/tunnel"

	"github.com/spf13/cobra"
)

// tcpCmd represents the http command
var tcpCmd = &cobra.Command{
	Use:   "tcp <port>",
	Short: "Expose a local tcp port you specify",
	Long: "Expose a local tcp port you specify.\n" +
		"Example: `punch tcp 8080` will expose a local tcp port running on port 8080.\n",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command) {
		var err error
		port = args[0]
		if err != nil {
			reportError("Must supply a port to forward", true)
		}
		tunnelTCP()
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)
}

func tunnelTCP() {
	if !checkPort(port) {
		reportError("Port is not in range[1-65535]", true)
	}

	publicKey, err := getPublicKey(publicKeyPath)
	if err != nil {
		os.Exit(3)
	}

	protocol := []string{"tcp"}
	response, err := restAPI.CreateTunnelAPI("", publicKey, protocol)

	if err != nil {
		reportError(err.Error(), true)
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

	tunnelConfig := tunnel.Config{
		ConnectionEndpoint: *connectionURL,
		RestAPI:            restAPI,
		TunnelEndpoint:     response,
		EndpointType:       "tcp",
		PrivateKeyPath:     privateKeyPath,
		EndpointURL:        *baseURL,
		LocalPort:          port,
		Subdomain:          "tcp",
		LogLevel:           logLevel,
	}
	semaphore := tunnel.Semaphore{}
	tunnel.StartReverseTunnel(&tunnelConfig, nil, &semaphore)
}
