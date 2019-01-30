package test

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestSubdomainReserve(t *testing.T) {
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	p = testcli.Command("../punch", "subdomain", "reserve", "testdomain", "--config", CONFIG_PATH)
	p.Run()

	require.Equal(t, p.Stdout(), "Successfully reserved subdomain testdomain\n")
	p = testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))

	defer func() {
		p := testcli.Command("../punch", "subdomain", "release", "testdomain", "--config", CONFIG_PATH)
		p.Run()
	}()
}
func TestOwnedSubdomainReserve(t *testing.T) {
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))

	p = testcli.Command("../punch", "subdomain", "reserve", "testdomain", "--config", CONFIG_PATH)
	p.Run()
	p = testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))

	p = testcli.Command("../punch", "subdomain", "reserve", "testdomain", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, p.Stdout(), "Subdomain has already been reserved\n")

	defer func() {
		p := testcli.Command("../punch", "subdomain", "release", "testdomain", "--config", CONFIG_PATH)
		p.Run()
	}()
}
