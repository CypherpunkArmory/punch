package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// reserveCmd represents the reserve command
var reserveCmd = &cobra.Command{
	Use:   "reserve [subdomain]",
	Short: "Reserve a subdomain",
	Long:  `Reserve a subdomain to secure the subdomain for future use. Once reserved only you can use it`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain = args[0]
		reserve()
	},
}

func init() {
	subdomainCmd.AddCommand(reserveCmd)

}

func reserve() {
	if !checkSubdomain(subdomain) {
		fmt.Fprintf(os.Stderr, "Invalid Subdomain\n")
		os.Exit(1)
	}

	response, err := restAPI.ReserveSubdomainAPI(subdomain)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	fmt.Println("Successfully reserved subdomain " + response.Name)
}
