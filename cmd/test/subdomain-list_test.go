package cmdtest

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestSubdomainListNoSubdomains(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsEmptySubdomainList(p.Stdout()))
}

func equalsEmptySubdomainList(output string) bool {
	return output == "+----------------+----------+--------+\n"+
		"| SUBDOMAIN NAME | RESERVED | IN USE |\n"+
		"+----------------+----------+--------+\n"+
		"+----------------+----------+--------+\n"
}
func TestSubdomainListWithSubdomains(t *testing.T) {
	defer createConfig(t)()
	configLogin(t)
	defer reserveSubdomain(t, "testdomain")()
	p := testcli.Command(exePath, "subdomain", "list", "--config", configPath)
	p.Run()
	require.Equal(t, true, equalsOneSubdomainList(p.Stdout(), "testdomain"))
}
func equalsOneSubdomainList(output string, subdomain string) bool {
	return output == "+----------------+----------+--------+\n"+
		"| SUBDOMAIN NAME | RESERVED | IN USE |\n"+
		"+----------------+----------+--------+\n"+
		"| "+subdomain+"     | true     | false  |\n"+
		"+----------------+----------+--------+\n"
}
