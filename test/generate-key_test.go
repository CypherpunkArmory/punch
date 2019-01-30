package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cypherpunkarmory/punch/utilities"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyNoParams(t *testing.T) {
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "generate-key", "--config", CONFIG_PATH)
	p.Run()
	keyPath := "../"
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
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "generate-key", "-n", "test_key", "--config", CONFIG_PATH)
	p.Run()
	keyPath := "../"
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
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "generate-key", "~/.ssh", "--config", CONFIG_PATH)
	p.Run()
	keyPath := utilities.FixFilePath("~/.ssh")
	pem, err := ioutil.ReadFile(keyPath + "/holepunch_key.pem")
	if err != nil {
		t.Fatal("Pem file not written")
	}
	pub, err := ioutil.ReadFile(keyPath + "/holepunch_key.pub")
	if err != nil {
		t.Fatal("Pub file not written")
	}
	defer func() {
		err := os.Remove(keyPath + "/holepunch_key.pem")
		if err != nil {
			t.Fatalf(keyPath + "holepunch_key.pem not deleted")
		}
		err = os.Remove(keyPath + "/holepunch_key.pub")
		if err != nil {
			t.Fatalf(keyPath + "holepunch_key.pub not deleted")
		}
	}()
	require.Contains(t, string(pem), "BEGIN RSA PRIVATE KEY")
	require.Contains(t, string(pub), "ssh-rsa ")
	require.Equal(t, p.Stdout(), "SSH keys have been generated and the config file has been updated\n")
}
