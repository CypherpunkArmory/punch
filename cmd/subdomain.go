package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cypherpunkarmory/punch/restapi"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var subdomainCmd = &cobra.Command{
	Use:   "subdomain",
	Short: "subdomain",
	Long:  `subdomain`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Try subdomain release, subdomain reserve or subdomain list")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your subdomains",
	Long:  `List subdomains you have previously reserved and also subdomains that are currently in use by you`,
	Run: func(cmd *cobra.Command, args []string) {
		subdomainList()
	},
}

func init() {
	rootCmd.AddCommand(subdomainCmd)
	subdomainCmd.AddCommand(listCmd)
}

func subdomainList() {
	response, err := restAPI.ListSubdomainAPI()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	printSubdomains(response)
}

func printSubdomains(response []restapi.Subdomain) {
	var data = make([][]string, len(response))
	for _, elem := range response {
		reserved := strconv.FormatBool(elem.Reserved)
		inuse := strconv.FormatBool(elem.InUse)
		data = append(data, []string{elem.Name, reserved, inuse})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subdomain Name", "Reserved", "In Use"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})

	table.AppendBulk(data)
	table.Render()
}
