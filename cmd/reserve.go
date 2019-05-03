package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// reserveCmd represents the reserve command
var reserveCmd = &cobra.Command{
	Use:   "reserve <subdomain>",
	Short: "Reserve a subdomain",
	Long: "Reserve a subdomain to secure the subdomain for future use.\n" +
		"Once reserved only you can use it.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain = args[0]
		reserve()
	},
}

func init() {
	rootCmd.AddCommand(reserveCmd)
}

func reserve() {
	if !correctSubdomainRegex(subdomain) {
		reportError("Invalid Subdomain", true)
	}

	response, err := restAPI.ReserveSubdomainAPI(subdomain)
	if err != nil {
		reportError(err.Error(), true)
	}
	fmt.Print("Successfully reserved subdomain " + response.Name)
	d := color.New(color.FgGreen, color.Bold)
	d.Printf(" ✔\n")
}
