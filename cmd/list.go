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
	"HolePunchCLI/restapi"
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your subdomains",
	Long:  `List subdomains you have previously reserved and also subdomains that are currently in use by you`,
	Run: func(cmd *cobra.Command, args []string) {
		restAPI := restapi.RestClient{
			URL:    BASE_URL,
			APIKEY: API_KEY,
		}
		response, err := restAPI.SubdomainListAPI()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			for _, elem := range response.Data {
				fmt.Printf("Name: %s\tReserved: %t\tInUse: %t\n", elem.Attributes.Name, elem.Attributes.Reserved, elem.Attributes.InUse)
			}
		}
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
