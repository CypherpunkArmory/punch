package cmd

import (
	"fmt"
	"os"

	"strconv"

	"github.com/cypherpunkarmory/punch/tunnel"
	"github.com/cypherpunkarmory/punch/utilities"
	"github.com/spf13/cobra"
)

// httpsCmd represents the https command
var httpsCmd = &cobra.Command{
	Use:   "https [port]",
	Short: "Expose a https web server over the port you specify",
	Long:  `To expose a https web server on port 443 punch https 443`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		port, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Must supply a port to forward")
			os.Exit(1)
		}
		tunnelHTTPS()
	},
}

func init() {
	rootCmd.AddCommand(httpsCmd)
	httpsCmd.Flags().StringVarP(&subdomain, "subdomain", "s", "", "If not selected domain name will be autogenerated")

}
func tunnelHTTPS() {
	if subdomain != "" && !utilities.CheckSubdomain(subdomain) {
		fmt.Println("Invalid Subdomain")
		os.Exit(1)
	}
	if !utilities.CheckPort(port) {
		fmt.Println("Port is not in range[1-65535")
		os.Exit(1)
	}

	publicKey, err := utilities.GetPublicKey(publicKeyPath)
	if err != nil {
		os.Exit(3)
	}

	protocol := []string{"https"}
	response, err := restAPI.CreateTunnelAPI(subdomain, publicKey, protocol)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if subdomain == "" {
		subdomain, _ = restAPI.GetSubdomainName(response.Subdomain.ID)
	}
	tunnelConfig := tunnel.Config{
		RestAPI:        restAPI,
		TunnelEndpoint: response,
		EndpointType:   "https",
		PrivateKeyPath: privateKeyPath,
		EndpointURL:    baseURL,
		LocalPort:      port,
		Subdomain:      subdomain,
	}
	tunnel.StartReverseTunnel(&tunnelConfig)

}
