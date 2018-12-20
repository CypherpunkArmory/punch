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
	"HolePunchCLI/utilities"
	"fmt"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release subdomain",
	Long:  `Release a subdomain you have reserved`,
	Run: func(cmd *cobra.Command, args []string) {
		if utilities.CheckSubdomain(Subdomain) {
			restAPI := restapi.RestClient{
				URL:    BASE_URL,
				APIKEY: API_KEY,
			}
			err := restAPI.ReleaseSubodmainAPI(Subdomain)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Successfully released subdomain")
			}
		} else {
			fmt.Println("Invalid Subdomain")
		}
	},
}

func init() {
	subdomainCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	releaseCmd.Flags().StringVarP(&Subdomain, "subdomain", "s", "", "")
}
