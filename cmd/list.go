package cmd

import (
	"fmt"
	"os"

	"text/tabwriter"

	"github.com/cypherpunkarmory/punch/restapi"
	"github.com/spf13/cobra"
)

var showTunnels bool
var showSubdomains bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your subdomains or tunnels",
	Long:  "List subdomains or tunnels you have previously reserved or that are currently in use by you.  By default, subdomains are listed",
	Run: func(cmd *cobra.Command, args []string) {
		if showTunnels {
			tunnelList()
		} else {
			subdomainList()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&showTunnels, "tunnels", "t", false, "List running tunnels")
	listCmd.Flags().BoolVarP(&showSubdomains, "subdomains", "s", false, "List owned subdomains")
}

func subdomainList() {
	response, err := restAPI.ListSubdomainAPI()
	if err != nil {
		reportError(err.Error(), true)
	}
	printSubdomains(response)
}

func printSubdomains(response []restapi.Subdomain) {
	if len(response) == 0 {
		fmt.Println("You have no subdomains")
		return
	}
	writer := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	writer.Init(os.Stdout, 16, 8, 0, '\t', 0)

	defer writer.Flush()
	fmt.Fprintf(writer, "%s\t%s\t%s\t", "Subdomain Name", "Reserved", "In Use")
	fmt.Fprintf(writer, "\n%s\t%s\t%s\t\n", "--------------", "--------", "------")
	for _, elem := range response {
		fmt.Fprintf(writer, "%s\t%t\t%t\t\n", elem.Name, elem.Reserved, elem.InUse)
	}
}

func tunnelList() {
	response, err := restAPI.ListTunnelsAPI()
	if err != nil {
		reportError(err.Error(), true)
	}
	printTunnels(response)
}

func printTunnels(response []restapi.Tunnel) {
	if len(response) == 0 {
		fmt.Println("You have no tunnels")
		return
	}
	writer := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	writer.Init(os.Stdout, 16, 8, 0, '\t', 0)

	defer writer.Flush()
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t\n", "Subdomain Name", "Port", "SSH Port", "IP Address")
	fmt.Fprintf(writer, "%s\t%s\t%s\t\n", "--------------", "--------", "------")
	for _, elem := range response {
		subdomainName, _ := restAPI.GetSubdomainName(elem.Subdomain.ID)
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t\n", subdomainName, elem.Port, elem.SSHPort, elem.IPAddress)
	}
}
