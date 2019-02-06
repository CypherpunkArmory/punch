package cmdtest

import (
	"github.com/cypherpunkarmory/punch/cmd"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	"testing"
)

var home, _ = homedir.Dir()

func TestRootTryStartSession(t *testing.T) {
	cmd.API_ENDPOINT = "http://0.0.0.0:5000"

	err := cmd.TryStartSession()
	require.NotNil(t, err)
}
