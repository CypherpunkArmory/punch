package test

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestSubdomainRelease(t *testing.T) {
	defer CreateConfig(t)()
	configLogin(t)
	//Make sure list is empty
	p := testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	//Make sure subdomain was reserved
	p = testcli.Command("../punch", "subdomain", "reserve", "testdomain", "--config", CONFIG_PATH)
	p.Run()
	p = testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))
	//Make sure subdomain was deleted
	p = testcli.Command("../punch", "subdomain", "release", "testdomain", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, p.Stdout(), "Successfully released subdomain\n")
	p = testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
}
func TestUnownedSubdomainRelease(t *testing.T) {
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	out := testcli.Command("../punch", "subdomain", "release", "testdomain", "--config", CONFIG_PATH)
	out.Run()

	require.Equal(t, out.Stdout(), "You do not own this subdomain\n")

}
