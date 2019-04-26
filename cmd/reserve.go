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
