package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup holepunch",
	Run: func(cmd *cobra.Command, args []string) {
		var setupKey string
		setupLogin()
		fmt.Print("Would you like to generate ssh keys to forward traffic? (Y/n): ")
		fmt.Scanln(&setupKey)
		setupKey = strings.ToLower(setupKey)
		if setupKey != "" && !strings.HasPrefix(setupKey, "y") && !strings.HasPrefix(setupKey, "n") {
			log.Println("Invalid input")
			os.Exit(1)
		}
		if strings.HasPrefix(setupKey, "n") {
			fmt.Println("Make sure you set the path to your keys in the config file located at: " + configPath +
				"\n You can also generate keys using the generate-key command")
			return
		}
		err := generateKey("", "holepunch_key")
		if err != nil {
			fmt.Println("Could not generate key")
			os.Exit(1)
		}
		fmt.Println("Generated keys in the current directory")
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
		fmt.Println("Error reading username")
		os.Exit(1)
	}
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error reading password")
		os.Exit(1)
	}
	fmt.Println()
	password = string(bytePassword)
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
