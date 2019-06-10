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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var clearAll bool

var cleanupCmd = &cobra.Command{
	Use:   "cleanup <subdomain>",
	Short: "Cleanup a subdomain that is incorrectly marked as \"In Use\"",
	Long: "Cleanup a subdomain that is incorrectly marked as \"In Use\".\n" +
		"This closes the tunnel from our end and updates the subdomain database.",
	Args: func(cmd *cobra.Command, args []string) error {
		if !clearAll && len(args) < 1 {
			return errors.New("requires a subdomain argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if clearAll {
			cleanupAll()
		} else {
			subdomain = args[0]
			cleanup(subdomain)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().BoolVarP(&clearAll, "all", "a", false, "Cleanup all subdomains incorrectly marked as in use")
}

func cleanup(openSubdomain string) {
	if !correctSubdomainRegex(openSubdomain) {
		reportError("Invalid Subdomain", true)
	}
	err := restAPI.DeleteTunnelAPI(openSubdomain)
	if err != nil {
		reportError(err.Error(), true)
	}
	fmt.Printf("%s%s\n", "Successfully closed tunnel for subdomain: ", openSubdomain)
}

func cleanupAll() {
	response, err := restAPI.ListSubdomainAPI()
	if err != nil {
		reportError(err.Error(), true)
	}

	for _, elem := range response {
		if !elem.InUse {
			continue
		}

		cleanup(elem.Name)
	}
}
