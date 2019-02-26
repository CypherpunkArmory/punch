package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var username string
var password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to holepunch",
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Your holepunch.io username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Your holepunch.io password")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}

func login() {
	response, err := restAPI.Login(username, password)

	if err != nil {
		fmt.Println("Login Failed: " + err.Error())
		os.Exit(1)
	}

	viper.Set("apikey", response.RefreshToken)
	err = viper.WriteConfig()

	if err != nil {
		fmt.Println("Couldn't write refresh token to config - permissions maybe?")
		os.Exit(1)
	}

}
