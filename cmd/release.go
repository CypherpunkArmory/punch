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

	"github.com/cypherpunkarmory/punch/utilities"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release [subdomain]",
	Short: "Release subdomain",
	Long:  `Release a subdomain you have reserved`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Subdomain = args[0]
		release(Subdomain)
	},
}

func init() {
	subdomainCmd.AddCommand(releaseCmd)
}

func release(Subdomain string) {
	if !utilities.CheckSubdomain(Subdomain) {
		fmt.Println("Invalid Subdomain")
		os.Exit(1)
	}
	err := restAPI.ReleaseSubdomainAPI(Subdomain)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully released subdomain")
}
