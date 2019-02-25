package cmdtest

import (
	"os"
	"runtime"
	"testing"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/rendon/testcli"
	"github.com/spf13/viper"
)

func getConfigPath() string {
	if runtime.GOOS == "windows" {
		home, _ := homedir.Dir()
		return home + "\\punch_test.toml"
	}
	return "/tmp/punch.toml"
}

func getExePath() string {
	if runtime.GOOS == "windows" {
		return ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "punch.exe"
	}
	return ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "punch"
}
func getKeyPath() string {
	if runtime.GOOS == "windows" {
		return ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator)
	}
	return ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator)
}

var configPath = getConfigPath()
var exePath = getExePath()

func createConfig(t *testing.T) func() {
	t.Helper()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.Create(configPath)
		initTestConfig(t)
	}
	return func() {
		err := os.Remove(configPath)
		if err != nil {
			viper.SetDefault("apikey", "")
			viper.SetDefault("baseurl", "holepunch.io")
			viper.SetDefault("apiendpoint", "http://0.0.0.0:5000")
			viper.SetDefault("publickeypath", getKeyPath())
			viper.SetDefault("privatekeypath", getKeyPath())
			err = viper.WriteConfigAs(configPath)
			if err != nil {
				t.Fatalf("Cant write config file")
			}
		}
	}
}
func initTestConfig(t *testing.T) {
	t.Helper()
	viper.SetDefault("apikey", "")
	viper.SetDefault("baseurl", "holepunch.io")
	viper.SetDefault("apiendpoint", "http://0.0.0.0:5000")
	viper.SetDefault("publickeypath", getKeyPath())
	viper.SetDefault("privatekeypath", getKeyPath())
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		t.Fatalf("Couldn't generate config file")
	}
}

func configLogin(t *testing.T) {
	t.Helper()
	p := testcli.Command(exePath, "login", "-u", "testuser@holepunch.io", "-p", "secret", "--config", configPath)
	p.Run()
}

func reserveSubdomain(t *testing.T, subdomain string) func() {
	t.Helper()
	p := testcli.Command(exePath, "subdomain", "reserve", subdomain, "--config", configPath)
	p.Run()
	return func() {
		p := testcli.Command(exePath, "subdomain", "release", subdomain, "--config", configPath)
		p.Run()
	}
}
