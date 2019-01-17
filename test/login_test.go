package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestAskForLogin(t *testing.T) {
	p := testcli.Command("../punch", "subdomain", "list")
	p.Run()

	require.Contains(t, p.Stdout(), "You need to login using `punch login` first.")
}

func TestLogin(t *testing.T) {
	p := testcli.Command("../punch", "login")
	p.Run()
	if !p.Failure() {
		t.Fatalf("Expected punch login to fail, but it succeeed.")
	}

	if !p.StdoutContains("required flag(s) \"password\", \"username\" not set") {
		t.Fatalf("Expected password and username to be required.")
	}
}

func TestLoginSetsTOML(t *testing.T) {
	p := testcli.Command("../punch", "login", "-u", "test@londontrustmedia.com", "-p", "test", "--baseurl", "http://localhost:5000")
	p.Run()

	if !p.Success() {
		t.Fatalf("Expected punch login to succeeed, but it failed.")
	}

	fmt.Println(p.Stdout())

	dat, err := ioutil.ReadFile("/tmp/punch.toml")
	if err != nil {
		t.Fatal("/tmp/punch.toml not written")
	}

	defer func() {
		err := os.Remove("/tmp/punch.toml")
		if err != nil {
			t.Fatalf("/tmp/punch.toml not deleted")
		}
	}()

	require.Contains(t, string(dat), "apikey=\"eyJ0eXAiO")

}
