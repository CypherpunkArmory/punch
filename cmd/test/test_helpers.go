package cmdtest

import (
	"fmt"
	"os"
	"testing"

	"github.com/rendon/testcli"
	"github.com/spf13/viper"
)

var CONFIG_PATH = "/tmp/punch.toml"

func createConfig(t *testing.T) func() {
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		os.Create(CONFIG_PATH)
		initConfig()
	}
	return func() {
		err := os.Remove(CONFIG_PATH)
		if err != nil {
			t.Fatalf(CONFIG_PATH + " not deleted")
		}
	}
}
func initConfig() {
	viper.SetDefault("apikey", "")
	viper.SetDefault("baseurl", "holepunch.io")
	viper.SetDefault("apiendpoint", "http://0.0.0.0:5000")
	viper.SetDefault("publickeypath", "~/.ssh/holepunch_key.pub")
	viper.SetDefault("privatekeypath", "~/.ssh/holepunch_key.pem")
	err := viper.WriteConfigAs(CONFIG_PATH)
	if err != nil {
		fmt.Println("Couldn't generate default config file")
	}
}

func configLogin(t *testing.T) {
	p := testcli.Command("../../punch", "login", "-u", "testuser@holepunch.io", "-p", "secret", "--config", CONFIG_PATH)
	p.Run()
}

func reserveSubdomain(t *testing.T, subdomain string) func() {
	p := testcli.Command("../../punch", "subdomain", "reserve", subdomain, "--config", CONFIG_PATH)
	p.Run()
	return func() {
		p := testcli.Command("../../punch", "subdomain", "release", subdomain, "--config", CONFIG_PATH)
		p.Run()
	}
}
