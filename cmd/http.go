package cmd

import (
	"net/url"
	"os"
	"strconv"

	"github.com/cypherpunkarmory/punch/tunnel"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http [port] [subdomain]",
	Short: "Expose a web server on the port you specify",
	Long:  `To expose a web server on port 80 punch http 80`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) == 2 {
			subdomain = args[1]
		}
		port, err = strconv.Atoi(args[0])
		if err != nil {
			reportError("Must supply a port to forward", true)
		}
		tunnelHTTP()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func tunnelHTTP() {
	if subdomain != "" && !checkSubdomain(subdomain) {
		reportError("Invalid Subdomain", true)
	}
	if !checkPort(port) {
		reportError("Port is not in range[1-65535]", true)
	}

	publicKey, err := getPublicKey(publicKeyPath)
	if err != nil {
		os.Exit(3)
	}

	protocol := []string{"http"}
	response, err := restAPI.CreateTunnelAPI(subdomain, publicKey, protocol)

	if err != nil {
		reportError(err.Error(), true)
	}

	if subdomain == "" {
		subdomain, _ = restAPI.GetSubdomainName(response.Subdomain.ID)
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
		EndpointType:       "http",
		PrivateKeyPath:     privateKeyPath,
		EndpointURL:        *baseURL,
		LocalPort:          port,
		Subdomain:          subdomain,
		LogLevel:           logLevel,
	}
	tunnel.StartReverseTunnel(&tunnelConfig, nil)
}
