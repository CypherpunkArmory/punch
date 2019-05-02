package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var username string
var password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to holepunch",
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
		reportError("Couldn't write refresh token to config - permissions maybe?", true)
	}
	fmt.Println("Login successful")
}
