package root_test

import (
	"github.com/cypherpunkarmory/punch/cmd"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	"os"
	"testing"
)

var home, _ = homedir.Dir()

func TestRootTryReadConfig(t *testing.T) {
	cmd.CfgFile = "/tmp/punch.toml"

	err := cmd.TryReadConfig()

	require.NotNil(t, err)
}

func TestRootTryReadConfigWithoutOverride(t *testing.T) {
	cmd.CfgFile = ""

	// make sure the default config doesn't exist
	_ = os.Remove(home + "/.punch.toml")

	err := cmd.TryReadConfig()

	require.Nil(t, err)

	stat, err := os.Stat(home + "/.punch.toml")

	require.NotNil(t, stat)
}

func TestRootTryStartSession(t *testing.T) {
	cmd.BASE_URL = "http://notathing.url"

	err := cmd.TryStartSession()

	require.NotNil(t, err)
}
