package test

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/require"
)

func TestSubdomainListNoSubdomains(t *testing.T) {
	defer CreateConfig(t)()
	configLogin(t)
	p := testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
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
	defer CreateConfig(t)()
	configLogin(t)
	defer reserveSubdomain(t, "testdomain")()
	p := testcli.Command("../punch", "subdomain", "list", "--config", CONFIG_PATH)
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
