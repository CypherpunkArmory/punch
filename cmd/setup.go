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
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup punch (run this first)",
	Long: "Setup punch.\n" +
		"This will ask you for your holepunch credentials and help you create pub/priv keys if needed.",
	Run: func(cmd *cobra.Command, args []string) {
		var setupKey string
		setupLogin()
		fmt.Print("Would you like to generate ssh keys to forward traffic? (Y/n): ")
		fmt.Scanln(&setupKey)
		setupKey = strings.ToLower(setupKey)
		if setupKey != "" && !strings.HasPrefix(setupKey, "y") && !strings.HasPrefix(setupKey, "n") {
			reportError("Invalid input", true)
		}
		if strings.HasPrefix(setupKey, "n") {
			fmt.Println("Make sure you set the path to your keys in the config file located at: " + configPath +
				"\n You can also generate keys using the generate-key command")
			return
		}
		err := generateKey("", "holepunch_key")
		if err != nil {
			reportError("Could not generate key", true)
		}
		fmt.Print("Generated keys in the current directory ")
		d := color.New(color.FgGreen, color.Bold)
		d.Printf("✔\n")
		err = writeKeysToConfig("", "holepunch_key")
		if err != nil {
			reportError("Failed to update config file", true)
		}
		fmt.Print("Config file updated")
		d.Printf(" ✔\n")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func setupLogin() {
	var username string
	var password string
	fmt.Print("Enter Username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		reportError("Error reading username", true)
	}
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		reportError("Error reading password", true)
	}
	fmt.Println()
	password = string(bytePassword)
	response, err := restAPI.Login(username, password)

	if err != nil {
		reportError("Login Failed: "+err.Error(), true)
	}

	viper.Set("apikey", response.RefreshToken)
	err = viper.WriteConfig()

	if err != nil {
		reportError("Couldn't write refresh token to config - permissions maybe?", true)
	}
	fmt.Print("Login Succesful ")
	d := color.New(color.FgGreen, color.Bold)
	d.Printf("✔\n")
}
