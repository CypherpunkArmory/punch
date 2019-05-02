package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup [subdomain]",
	Short: "Cleanup a subdomain that is incorrectly marked as \"In Use\"",
	Long: "Cleanup a subdomain that is incorrectly marked as \"In Use\".\n" +
		"This closes the tunnel from our end and updates the subdomain database.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain = args[0]
		cleanup(subdomain)
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}

func cleanup(openSubdomain string) {
	if !correctSubdomainRegex(openSubdomain) {
		reportError("Invalid Subdomain", true)
	}
	err := restAPI.DeleteTunnelAPI(openSubdomain)
	if err != nil {
		reportError(err.Error(), true)
	}
	fmt.Println("Successfully closed tunnel")

}
