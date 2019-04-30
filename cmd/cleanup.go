package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup [subdomain]",
	Short: "Cleanup a tunnel that shouldn't be open",
	Long:  `Cleanup a tunnel that shouldn't be open that is associated to the given subdomain`,
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
		reportError("Invalid Subdomain", true)
	}
	err := restAPI.DeleteTunnelAPI(openSubdomain)
	if err != nil {
		reportError(err.Error(), true)
	}
	fmt.Println("Successfully closed tunnel")

}
