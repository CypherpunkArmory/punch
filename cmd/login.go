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
	"github.com/spf13/viper"
)

var username string
var password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to holepunch.",
	Long: "Login back into holepunch.\n" +
		"Will prompt you for username and password, or you can provide them as optional arguments.\n" +
		"If this is your first time using punch, you should use `punch setup` instead of `punch login`.",
	Run: func(cmd *cobra.Command, args []string) {
		if username != "" && password != "" {
			login()
			return
		}
		setupLogin() // This function is located in cmd/setup.go
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Your holepunch.io username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Your holepunch.io password")
}

func login() {
	response, err := restAPI.Login(username, password)

	if err != nil {
		reportError("Login Failed: "+err.Error(), true)
	}

	viper.Set("apikey", response.RefreshToken)
	err = viper.WriteConfig()

	if err != nil {
		reportError("Couldn't write refresh token to config - are you able to write to ~/.config/holepunch/punch.toml?", true)
	}
	fmt.Print("Login Succesful ")
	d := color.New(color.FgGreen, color.Bold)
	d.Printf("✔\n")
}
