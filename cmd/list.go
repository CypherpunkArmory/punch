// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cypherpunkarmory/punch/restapi"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your subdomains",
	Long:  `List subdomains you have previously reserved and also subdomains that are currently in use by you`,
	Run: func(cmd *cobra.Command, args []string) {
		subdomainList()
	},
}

func init() {
	subdomainCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func subdomainList() {
	response, err := restAPI.ListSubdomainAPI()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		printSubdomains(&response)
	}
}

func printSubdomains(response *[]restapi.Subdomain) {
	var data [][]string
	for _, elem := range *response {
		reserved := strconv.FormatBool(elem.Reserved)
		inuse := strconv.FormatBool(elem.InUse)
		data = append(data, []string{elem.Name, reserved, inuse})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Subdomain Name", "Reserved", "In Use"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})

	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})

	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor})

	table.AppendBulk(data)
	table.Render()
}
