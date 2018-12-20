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
	"HolePunchCLI/restapi"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Port int
var Subdomain string
var Verbose bool
var REFRESH_TOKEN string
var API_KEY string
var BASE_URL string
var PUBLIC_KEY_PATH string
var PRIVATE_KEY_PATH string

var rootCmd = &cobra.Command{
	Use:   "punch",
	Short: "Like a holepunch for your network",
	Long:  `HolePunch`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//var Verbose bool
	cobra.OnInitialize(initConfig)

	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
	viper.BindPFlag("publickeypath", rootCmd.PersistentFlags().Lookup("publickeypath"))
	viper.BindPFlag("privatekeypath", rootCmd.PersistentFlags().Lookup("privatekeypath"))
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".HolePunch" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".HolePunch")
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		REFRESH_TOKEN = viper.GetString("REFRESH_TOKEN")
		BASE_URL = viper.GetString("BASE_URL")
		PUBLIC_KEY_PATH = viper.GetString("PUBLIC_KEY_PATH")
		PRIVATE_KEY_PATH = viper.GetString("PRIVATE_KEY_PATH")
		restAPI := restapi.RestClient{
			URL: BASE_URL,
		}
		response, err := restAPI.StartSession(REFRESH_TOKEN)
		if err != nil {

		}
		API_KEY = response.Access_Token
	} else {
		fmt.Println(err)
	}

}
