package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release -s [subdomain]",
	Short: "Release subdomain",
	Long:  `Release a subdomain you have reserved`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain = args[0]
		release(subdomain)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func release(subdomain string) {
	if !correctSubdomainRegex(subdomain) {
		reportError("Invalid Subdomain", true)
	}
	err := restAPI.ReleaseSubdomainAPI(subdomain)
	if err != nil {
		reportError(err.Error(), true)
	}
	fmt.Print("Successfully released subdomain ")
	d := color.New(color.FgGreen, color.Bold)
	d.Printf("✔\n")
}
