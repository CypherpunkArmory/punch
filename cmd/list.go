package cmd

import (
	"fmt"
	"os"

	"text/tabwriter"

	"github.com/cypherpunkarmory/punch/restapi"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your subdomains",
	Long:  `List subdomains you have previously reserved and also subdomains that are currently in use by you`,
	Run: func(cmd *cobra.Command, args []string) {
		subdomainList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func subdomainList() {
	response, err := restAPI.ListSubdomainAPI()
	if err != nil {
		reportError(err.Error(), true)
	}
	printSubdomains(response)
}

func printSubdomains(response []restapi.Subdomain) {
	w := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 16, 8, 0, '\t', 0)

	defer w.Flush()
	if len(response) == 0 {
		fmt.Println("You have no subdomains")
		return
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t", "Subdomain Name", "Reserved", "In Use")
	fmt.Fprintf(w, "\n%s\t%s\t%s\t\n", "--------------", "--------", "------")
	for _, elem := range response {
		fmt.Fprintf(w, "%s\t%t\t%t\t\n", elem.Name, elem.Reserved, elem.InUse)
	}
}