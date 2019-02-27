// +build all integration

package cmdtest

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyNoParams(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	p := testcli.Command(exePath, "generate-key", "--config", configPath)
	p.Run()
	keyPath := getKeyPath()
	pem, err := ioutil.ReadFile(keyPath + "holepunch_key.pem")
	if err != nil {
		t.Fatal("Pem file not written" + keyPath)
	}
	pub, err := ioutil.ReadFile(keyPath + "holepunch_key.pub")
	if err != nil {
		t.Fatal("Pub file not written")
	}
	defer func() {
		err := os.Remove(keyPath + "holepunch_key.pem")
		if err != nil {
			t.Fatalf(keyPath + "holepunch_key.pem not deleted")
		}
		err = os.Remove(keyPath + "holepunch_key.pub")
		if err != nil {
			t.Fatalf(keyPath + "holepunch_key.pub not deleted")
		}
	}()
	require.Contains(t, string(pem), "BEGIN RSA PRIVATE KEY")
	require.Contains(t, string(pub), "ssh-rsa ")
	require.Equal(t, p.Stdout(), "SSH keys have been generated and the config file has been updated\n")
}
func TestGenerateKeyWithName(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	p := testcli.Command(exePath, "generate-key", "-n", "test_key", "--config", configPath)
	p.Run()
	keyPath := getKeyPath()
	pem, err := ioutil.ReadFile(keyPath + "test_key.pem")
	if err != nil {
		t.Fatal("Pem file not written" + keyPath)
	}
	pub, err := ioutil.ReadFile(keyPath + "test_key.pub")
	if err != nil {
		t.Fatal("Pub file not written")
	}
	defer func() {
		err := os.Remove(keyPath + "test_key.pem")
		if err != nil {
			t.Fatalf(keyPath + "test_key.pem not deleted")
		}
		err = os.Remove(keyPath + "test_key.pub")
		if err != nil {
			t.Fatalf(keyPath + "test_key.pub not deleted")
		}
	}()
	require.Contains(t, string(pem), "BEGIN RSA PRIVATE KEY")
	require.Contains(t, string(pub), "ssh-rsa ")
	require.Equal(t, p.Stdout(), "SSH keys have been generated and the config file has been updated\n")
}
func TestGenerateKeyWithLocation(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	path := "/tmp"
	if runtime.GOOS == "windows" {
		home, _ := homedir.Dir()
		path = home + "\\AppData\\Local\\Temp"
	}
	p := testcli.Command(exePath, "generate-key", path, "--config", configPath)
	p.Run()
	pem, err := ioutil.ReadFile(path + "/holepunch_key.pem")
	if err != nil {
		t.Fatal("Pem file not written")
	}
	pub, err := ioutil.ReadFile(path + "/holepunch_key.pub")
	if err != nil {
		t.Fatal("Pub file not written")
	}
	defer func() {
		err := os.Remove(path + "/holepunch_key.pem")
		if err != nil {
			t.Fatalf(path + "holepunch_key.pem not deleted")
		}
		err = os.Remove(path + "/holepunch_key.pub")
		if err != nil {
			t.Fatalf(path + "holepunch_key.pub not deleted")
		}
	}()
	require.Contains(t, string(pem), "BEGIN RSA PRIVATE KEY")
	require.Contains(t, string(pub), "ssh-rsa ")
	require.Equal(t, p.Stdout(), "SSH keys have been generated and the config file has been updated\n")
}
