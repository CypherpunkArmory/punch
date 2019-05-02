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
	Long:  "Login back into holepunch.\n" +
	       "You should use `punch setup` instead of `punch login` the first time.",
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
