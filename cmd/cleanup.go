package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup [subdomain]",
	Short: "cleanup a tunnel that shouldn't be open",
	Long:  `cleanup a tunnel that shouldn't be open that is associated to the given subdomain`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain = args[0]
		cleanup(subdomain)
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}

func cleanup(openSubdomain string) {
	if !checkSubdomain(openSubdomain) {
		fmt.Println("Invalid Subdomain")
		os.Exit(1)
	}
	err := restAPI.DeleteTunnelAPI(openSubdomain)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	fmt.Println("Successfully closed tunnel")

}
