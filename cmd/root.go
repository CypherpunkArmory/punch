// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cypherpunkarmory/punch/restapi"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string
var Port int
var Subdomain string
var Verbose bool
var REFRESH_TOKEN string
var API_KEY string
var API_ENDPOINT string
var PUBLIC_KEY_PATH string
var PRIVATE_KEY_PATH string
var BASE_URL string

var restAPI restapi.RestClient

var rootCmd = &cobra.Command{
	Version: "v0.2",
	Use:     "punch",
	Short:   "Like a holepunch for your network",
	Long:    `HolePunch`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		TryStartSession()
	},
}

// I can't imagine a situation in which this fails - non login shells?
var home, _ = homedir.Dir()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is ~/.punch)")
	rootCmd.PersistentFlags().StringVar(&API_KEY, "apikey", "", "Your holepunch API key")
	rootCmd.PersistentFlags().StringVar(&BASE_URL, "baseurl", "", "Holepunch server to use - (default is holepunch.io)")
	rootCmd.PersistentFlags().StringVar(&API_ENDPOINT, "apiendpoint", "", "Holepunch server to use - (default is http://api.holepunch.io)")
	rootCmd.PersistentFlags().StringVar(&PUBLIC_KEY_PATH, "publickeypath", "", "Path to your public keys - (~/.ssh)")
	rootCmd.PersistentFlags().StringVar(&PRIVATE_KEY_PATH, "privatekeypath", "", "Path to your private keys - (~/.ssh)")

	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
	viper.BindPFlag("apiendpoint", rootCmd.PersistentFlags().Lookup("apiendpoint"))
	viper.BindPFlag("publickeypath", rootCmd.PersistentFlags().Lookup("publickeypath"))
	viper.BindPFlag("privatekeypath", rootCmd.PersistentFlags().Lookup("privatekeypath"))
	viper.SetDefault("baseurl", "holepunch.io")
	viper.SetDefault("apiendpoint", "http://api.holepunch.io")
	viper.SetDefault("publickeypath", "~/.ssh/holepunch_key.pub")
	viper.SetDefault("privatekeypath", "~/.ssh/holepunch_key.pem")
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")
	viper.AddConfigPath(home)
	viper.AddConfigPath(home + "/.config/holepunch/")
	viper.SetConfigName(".punch")

	err := TryReadConfig()
	if err != nil {
		os.Exit(1)
	}

	viper.AutomaticEnv() // read in environment variables that match
}

func TryStartSession() {
	if REFRESH_TOKEN == "" {
		fmt.Println("You need to login using `punch login` first.")
		os.Exit(1)
	}

	restAPI = restapi.RestClient{
		URL:          API_ENDPOINT,
		RefreshToken: REFRESH_TOKEN,
	}

	// StartSession will set the internal state of the RestClient
	// to the correct API key
	_, err := restAPI.StartSession(REFRESH_TOKEN)

	if err != nil {
		fmt.Println("Error starting session")
		fmt.Println("You need to login using `punch login` first.")
		os.Exit(1)
	}
}

func TryReadConfig() (err error) {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
		if _, err := os.Stat(CfgFile); os.IsNotExist(err) {
			fmt.Println("Config file does not exist.")
			return err
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		REFRESH_TOKEN = viper.GetString("apikey")
		BASE_URL = viper.GetString("baseurl")
		PUBLIC_KEY_PATH = viper.GetString("publickeypath")
		PRIVATE_KEY_PATH = viper.GetString("privatekeypath")
		API_ENDPOINT = viper.GetString("apiendpoint")
	} else {
		if _, err := os.Stat(home + "/.punch.toml"); err != nil {
			if os.IsNotExist(err) {
				err := viper.WriteConfigAs(home + "/.punch.toml")
				if err != nil {
					fmt.Println("Couldn't generate default config file")
					return err
				}
			}
		} else {
			fmt.Println("You have an issue in your current config")
			return errors.New("Configuration Error")
		}

		fmt.Println("Generated default config.")
		_ = TryReadConfig()
	}

	return nil
}
