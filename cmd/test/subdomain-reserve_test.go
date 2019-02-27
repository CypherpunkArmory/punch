// +build all integration

package cmdtest

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestBadSubdomainReserve(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	p = testcli.Command(exePath, "subdomain", "reserve", "testdomain*/*/*/$$", "--config", configPath)
	p.Run()

	require.Equal(t, p.Stdout(), "Invalid Subdomain\n")
	p = testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
}
func TestSubdomainReserve(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)

	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
	p = testcli.Command(exePath, "subdomain", "reserve", "testdomain", "--config", configPath)
	p.Run()

	require.Equal(t, p.Stdout(), "Successfully reserved subdomain testdomain\n")
	p = testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))
	defer func() {
		p := testcli.Command(exePath, "subdomain", "release", "testdomain", "--config", configPath)
		p.Run()
	}()

}
func TestOwnedSubdomainReserve(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))

	p = testcli.Command(exePath, "subdomain", "reserve", "testdomain", "--config", configPath)
	p.Run()
	p = testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))

	p = testcli.Command(exePath, "subdomain", "reserve", "testdomain", "--config", configPath)
	p.Run()
	require.Equal(t, p.Stdout(), "Subdomain has already been reserved\n")

	defer func() {
		p := testcli.Command(exePath, "subdomain", "release", "testdomain", "--config", configPath)
		p.Run()
	}()
}
