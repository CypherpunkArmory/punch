// +build all integration

package cmdtest

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestSubdomainRelease(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	//Make sure list is empty
	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	//Make sure subdomain was reserved
	p = testcli.Command(exePath, "subdomain", "reserve", "testdomain", "--config", configPath)
	p.Run()
	p = testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))
	//Make sure subdomain was deleted
	p = testcli.Command(exePath, "subdomain", "release", "testdomain", "--config", configPath)
	p.Run()
	require.Equal(t, p.Stdout(), "Successfully released subdomain\n")
	p = testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
}
func TestUnownedSubdomainRelease(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	out := testcli.Command(exePath, "subdomain", "release", "testdomain", "--config", configPath)
	out.Run()

	require.Equal(t, out.Stderr(), "you do not own this subdomain\n")

}
