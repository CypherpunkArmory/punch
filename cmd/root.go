package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cypherpunkarmory/punch/restapi"
	rollbar "github.com/rollbar/rollbar-go"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiEndpoint string
var apiToken string
var baseURL string
var configFile string
var configPath string
var crashReporting bool
var port int
var privateKeyPath string
var publicKeyPath string
var refreshToken string
var restAPI restapi.RestClient
var rollbarToken string
var subdomain string

//This gets written in the makefile
var version string

var rootCmd = &cobra.Command{
	Version: version,
	Use:     "punch",
	Short:   "Like a holepunch for your network",
	Long:    `HolePunch`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		err := tryStartSession()
		if err != nil {
			os.Exit(1)
		}
		rollbar.SetToken(rollbarToken)
		rollbar.SetEnvironment("production")
		rollbar.SetCodeVersion(version)
		rollbar.SetServerRoot("github.com/CypherpunkArmory/punch")
		rollbar.SetCaptureIp(rollbar.CaptureIpNone)
		rollbar.SetEnabled(crashReporting && apiEndpoint == "http://api.holepunch.io")
	},
}

// I can't imagine a situation in which this fails - non login shells?
var home, _ = homedir.Dir()

//Execute is the entrypoint of cmd calls
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.punch)")
	rootCmd.PersistentFlags().StringVar(&apiToken, "apikey", "", "Your holepunch API key")
	rootCmd.PersistentFlags().StringVar(&baseURL, "baseurl", "", "Holepunch server to use - (default is holepunch.io)")
	rootCmd.PersistentFlags().StringVar(&apiEndpoint, "apiendpoint", "", "Holepunch server to use - (default is http://api.holepunch.io)")
	rootCmd.PersistentFlags().StringVar(&publicKeyPath, "publickeypath", "", "Path to your public keys - (~/.ssh)")
	rootCmd.PersistentFlags().StringVar(&privateKeyPath, "privatekeypath", "", "Path to your private keys - (~/.ssh)")
	rootCmd.PersistentFlags().BoolVar(&crashReporting, "crashreporting", false, "Send crash reports to the developers")

	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
	viper.BindPFlag("apiendpoint", rootCmd.PersistentFlags().Lookup("apiendpoint"))
	viper.BindPFlag("publickeypath", rootCmd.PersistentFlags().Lookup("publickeypath"))
	viper.BindPFlag("privatekeypath", rootCmd.PersistentFlags().Lookup("privatekeypath"))
	viper.BindPFlag("crashreporting", rootCmd.PersistentFlags().Lookup("crashreporting"))
	viper.SetDefault("crashreporting", false)
	viper.SetDefault("baseurl", "holepunch.io")
	viper.SetDefault("apiendpoint", "http://api.holepunch.io")
	viper.SetDefault("publickeypath", "")
	viper.SetDefault("privatekeypath", "")
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")
	//https://0x46.net/thoughts/2019/02/01/dotfile-madness/
	configPath = os.Getenv("XDG_CONFIG_HOME")
	if configPath == "" {
		configPath = home
	}
	configPath = filepath.Join(configPath, ".config", "holepunch")
	viper.AddConfigPath(configPath)
	viper.SetConfigName(".punch")

	err := tryReadConfig()
	if err != nil {
		os.Exit(1)
	}

	viper.AutomaticEnv() // read in environment variables that match
}

func tryStartSession() error {
	if refreshToken == "" {
		fmt.Println("You need to login using `punch login` first.")
		return errors.New("no refresh token")
	}

	// StartSession will set the internal state of the RestClient
	// to the correct API key
	err := restAPI.StartSession(refreshToken)

	if err != nil {
		fmt.Println("Error starting session")
		fmt.Println("You need to login using `punch login` first.")
		return errors.New("error starting session")
	}
	return nil
}

func tryReadConfig() (err error) {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Println("Config file does not exist.")
			return err
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		refreshToken = viper.GetString("apikey")
		baseURL = viper.GetString("baseurl")
		publicKeyPath = viper.GetString("publickeypath")
		privateKeyPath = viper.GetString("privatekeypath")
		apiEndpoint = viper.GetString("apiendpoint")

		publicKeyPath = fixFilePath(publicKeyPath)
		privateKeyPath = fixFilePath(privateKeyPath)
	} else {
		if _, err := os.Stat(configPath + string(os.PathSeparator) + ".punch.toml"); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(configPath, os.ModePerm)
				if err != nil {
					fmt.Println("Couldn't generate default config file")
					return err
				}
				err = viper.WriteConfigAs(configPath + string(os.PathSeparator) + ".punch.toml")
				if err != nil {
					fmt.Println("Couldn't generate default config file")
					return err
				}
			}
		} else {
			fmt.Println("You have an issue in your current config")
			return errors.New("configuration error")
		}

		fmt.Println("Generated default config.")
		_ = tryReadConfig()
	}
	restAPI = restapi.RestClient{
		URL:          apiEndpoint,
		RefreshToken: refreshToken,
	}
	return nil
}
