package cmd

import (
	"net/url"
	"os"

	"github.com/cypherpunkarmory/punch/tunnel"
	"github.com/spf13/cobra"
)

// httpsCmd represents the https command
var httpsCmd = &cobra.Command{
	Use:   "https [port]",
	Short: "Expose a https web server over the port you specify",
	Long:  `To expose a https web server on port 443 punch https 443`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) == 2 {
			subdomain = args[1]
		}
		port = args[0]
		if err != nil {
			reportError("Must supply a port to forward", true)
		}
		tunnelHTTPS()
	},
}

func init() {
	rootCmd.AddCommand(httpsCmd)

}
func tunnelHTTPS() {
	if subdomain != "" && !correctSubdomainRegex(subdomain) {
		reportError("Invalid subdomain", true)
	}
	if !checkPort(port) {
		reportError("Port is not in range[1-65535]", true)
	}

	publicKey, err := getPublicKey(publicKeyPath)
	if err != nil {
		os.Exit(3)
	}

	protocol := []string{"https"}
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
		EndpointType:       "https",
		PrivateKeyPath:     privateKeyPath,
		EndpointURL:        *baseURL,
		LocalPort:          port,
		Subdomain:          subdomain,
		LogLevel:           logLevel,
	}
	semaphore := tunnel.Semaphore{}
	tunnel.StartReverseTunnel(&tunnelConfig, nil, &semaphore)
}
