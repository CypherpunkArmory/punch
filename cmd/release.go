// Punch CLI used for interacting with holepunch.io
// Copyright (C) 2018-2019  Orb.House, LLC
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release <subdomain>",
	Short: "Release subdomain",
	Long:  "Release a subdomain you have reserved.",
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
