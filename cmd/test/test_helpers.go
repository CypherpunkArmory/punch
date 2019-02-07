package cmdtest

import (
	"os"
	"testing"

	"github.com/rendon/testcli"
	"github.com/spf13/viper"
)

var CONFIG_PATH = "/tmp/punch.toml"

func createConfig(t *testing.T) func() {
	t.Helper()
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		os.Create(CONFIG_PATH)
		initTestConfig(t)
	}
	return func() {
		err := os.Remove(CONFIG_PATH)
		if err != nil {
			t.Fatalf(CONFIG_PATH + " not deleted")
		}
	}
}
func initTestConfig(t *testing.T) {
	t.Helper()
	viper.SetDefault("apikey", "")
	viper.SetDefault("baseurl", "holepunch.io")
	viper.SetDefault("apiendpoint", "http://0.0.0.0:5000")
	viper.SetDefault("publickeypath", "/tmp/holepunch_key.pub")
	viper.SetDefault("privatekeypath", "/tmp/holepunch_key.pem")
	err := viper.WriteConfigAs(CONFIG_PATH)
	if err != nil {
		t.Fatalf("Couldn't generate config file")
	}
}

func configLogin(t *testing.T) {
	t.Helper()
	p := testcli.Command("../../punch", "login", "-u", "testuser@holepunch.io", "-p", "secret", "--config", CONFIG_PATH)
	p.Run()
}

func reserveSubdomain(t *testing.T, subdomain string) func() {
	t.Helper()
	p := testcli.Command("../../punch", "subdomain", "reserve", subdomain, "--config", CONFIG_PATH)
	p.Run()
	return func() {
		p := testcli.Command("../../punch", "subdomain", "release", subdomain, "--config", CONFIG_PATH)
		p.Run()
	}
}
