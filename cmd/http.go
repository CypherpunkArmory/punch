package cmd

import (
	"net/url"
	"os"

	"github.com/cypherpunkarmory/punch/tunnel"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http <port> [subdomain]",
	Short: "Expose a local web server on the port you specify",
	Long: "Expose a local web server on the port you specify.\n" +
		"Example: `punch http 8080` will expose a local web server running on port 8080.\n" +
		"You can provide an optional 2nd argument to specify the name of a reserved subdomain you want to\n" +
		"associate this with.\n" +
		"Example: `punch http 8080 mydomain` will expose a local web server running on port 8080 via\n" +
		"         \"http://mydomain.holepunch.io\".\n" +
		"Otherwise it will default to using a new unreserved subdomain.",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) == 2 {
			subdomain = args[1]
		}
		port = args[0]
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
	if subdomain != "" && !correctSubdomainRegex(subdomain) {
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
	semaphore := tunnel.Semaphore{}
	tunnel.StartReverseTunnel(&tunnelConfig, nil, &semaphore)
}
