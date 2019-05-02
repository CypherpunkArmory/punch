package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
var port string
var privateKeyPath string
var publicKeyPath string
var refreshToken string
var restAPI restapi.RestClient
var rollbarToken string
var sshEndpoint string
var subdomain string
var logLevel string

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
		rollbar.SetEnabled(crashReporting && apiEndpoint == "https://api.holepunch.io")
	},
}

// I can't imagine a situation in which this fails - non login shells?
var home, _ = homedir.Dir()

//Execute is the entrypoint of cmd calls
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		reportError(err.Error(), true)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.punch)")
	rootCmd.PersistentFlags().StringVar(&apiToken, "apikey", "", "Your holepunch API key")
	rootCmd.PersistentFlags().StringVar(&baseURL, "baseurl", "", "Holepunch server to use - (default is holepunch.io)")
	rootCmd.PersistentFlags().StringVar(&apiEndpoint, "apiendpoint", "", "Holepunch server to use - (default is https://api.holepunch.io)")
	rootCmd.PersistentFlags().StringVar(&publicKeyPath, "publickeypath", "", "Path to your public keys - (~/.ssh)")
	rootCmd.PersistentFlags().StringVar(&privateKeyPath, "privatekeypath", "", "Path to your private keys - (~/.ssh)")
	rootCmd.PersistentFlags().StringVar(&sshEndpoint, "sshendpoint", "", "endpoint that ssh tunnel connects to, to get to internal network")
	rootCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "", "Set the loglevel")
	rootCmd.PersistentFlags().BoolVar(&crashReporting, "crashreporting", false, "Send crash reports to the developers")
	err := rootCmd.PersistentFlags().MarkHidden("loglevel")
	if err != nil {
		panic(err)
	}

	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
	viper.BindPFlag("apiendpoint", rootCmd.PersistentFlags().Lookup("apiendpoint"))
	viper.BindPFlag("publickeypath", rootCmd.PersistentFlags().Lookup("publickeypath"))
	viper.BindPFlag("privatekeypath", rootCmd.PersistentFlags().Lookup("privatekeypath"))
	viper.BindPFlag("crashreporting", rootCmd.PersistentFlags().Lookup("crashreporting"))
	viper.BindPFlag("sshendpoint", rootCmd.PersistentFlags().Lookup("sshendpoint"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.SetDefault("crashreporting", true)
	viper.SetDefault("baseurl", "http://holepunch.io")
	viper.SetDefault("sshendpoint", "ssh://api.holepunch.io:22")
	viper.SetDefault("apiendpoint", "https://api.holepunch.io")
	viper.SetDefault("publickeypath", "")
	viper.SetDefault("privatekeypath", "")
	viper.SetDefault("loglevel", "ERROR")
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
		reportError("You need to login using `punch login` first.", false)
		return errors.New("no refresh token")
	}

	// StartSession will set the internal state of the RestClient
	// to the correct API key
	err := restAPI.StartSession(refreshToken)

	if err != nil {
		if strings.Contains(err.Error(), "incompatiable with the api") {
			reportError("Your punch client is out of date. Please use `punch update` to get the latest version.", false)
			confirmAndSelfUpdate()
			return errors.New("error starting session")
		}
		reportError("Error starting session", false)
		reportError("You need to login using `punch login` first.", false)
		return errors.New("error starting session")
	}
	return nil
}

func tryReadConfig() (err error) {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			reportError("Config file does not exist.", false)
			return err
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		refreshToken = viper.GetString("apikey")
		baseURL = viper.GetString("baseurl")
		publicKeyPath = viper.GetString("publickeypath")
		privateKeyPath = viper.GetString("privatekeypath")
		apiEndpoint = viper.GetString("apiendpoint")
		sshEndpoint = viper.GetString("sshendpoint")
		crashReporting = viper.GetBool("crashreporting")
		logLevel = viper.GetString("loglevel")

		publicKeyPath = fixFilePath(publicKeyPath)
		privateKeyPath = fixFilePath(privateKeyPath)
	} else {
		if _, err := os.Stat(configPath + string(os.PathSeparator) + ".punch.toml"); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(configPath, os.ModePerm)
				if err != nil {
					reportError("Couldn't generate default config file", false)
					return err
				}
				err = viper.WriteConfigAs(configPath + string(os.PathSeparator) + ".punch.toml")
				if err != nil {
					reportError("Couldn't generate default config file", false)
					return err
				}
			}
		} else {
			reportError("You have an issue in your current config", false)
			return errors.New("configuration error")
		}

		fmt.Println("Generated default config.")
		_ = tryReadConfig()
	}
	restAPI = restapi.NewRestClient(apiEndpoint, refreshToken)
	return nil
}
