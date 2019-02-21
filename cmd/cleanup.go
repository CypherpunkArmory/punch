package cmd

import (
	"fmt"
	"os"

	"github.com/cypherpunkarmory/punch/utilities"

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

func cleanup(Subdomain string) {
	if !utilities.CheckSubdomain(Subdomain) {
		fmt.Println("Invalid Subdomain")
		os.Exit(1)
	}
	err := restAPI.DeleteTunnelAPI(Subdomain)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully closed tunnel")

}
