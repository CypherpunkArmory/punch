package cmdtest

import (
	"os"
	"testing"

	"github.com/rendon/testcli"
	"github.com/spf13/viper"
)

var configPath = "/tmp/punch.toml"

func createConfig(t *testing.T) func() {
	t.Helper()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.Create(configPath)
		initTestConfig(t)
	}
	return func() {
		err := os.Remove(configPath)
		if err != nil {
			t.Fatalf(configPath + " not deleted")
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
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		t.Fatalf("Couldn't generate config file")
	}
}

func configLogin(t *testing.T) {
	t.Helper()
	p := testcli.Command("../../punch", "login", "-u", "testuser@holepunch.io", "-p", "secret", "--config", configPath)
	p.Run()
}

func reserveSubdomain(t *testing.T, subdomain string) func() {
	t.Helper()
	p := testcli.Command("../../punch", "subdomain", "reserve", subdomain, "--config", configPath)
	p.Run()
	return func() {
		p := testcli.Command("../../punch", "subdomain", "release", subdomain, "--config", configPath)
		p.Run()
	}
}
