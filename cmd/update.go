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
	"log"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update CLI version",
	Long:  "Update CLI to latest release on github.",
	Run: func(cmd *cobra.Command, args []string) {
		confirmAndSelfUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func confirmAndSelfUpdate() {
	up, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		log.Println("Couldn't create updater:", err)
		return
	}
	latest, found, err := up.DetectLatest("CypherpunkArmory/punch")
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return
	}
	// If version is not set just assume it should be updated
	v, _ := semver.Parse(version)
	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return
	}
	var input string
	fmt.Print("Do you want to update to: ", latest.Version, "? (Y/n): ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input != "" && !strings.HasPrefix(input, "y") && !strings.HasPrefix(input, "n") {
		log.Println("Invalid input")
		return
	}
	if strings.HasPrefix(input, "n") {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return
	}
	if err := up.UpdateTo(latest, exe); err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
}
