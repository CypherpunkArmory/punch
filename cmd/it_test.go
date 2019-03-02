// +build all unit

package cmd

import (
	"testing"
)

func Test_getTunnelConfig(t *testing.T) {
	cases := []struct {
		Name        string
		Input       string
		Expected    tunnelConf
		shouldError bool
	}{
		{"Valid", "http:80", tunnelConf{80, "http"}, false},
		{"Invalid(gibberish)", "kajsbdf&*(&", tunnelConf{}, true},
		{"Invalid(wrong order)", "80:http", tunnelConf{0, "80"}, true},
		{"Invalid(Extra values)", "http:80:https", tunnelConf{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := getTunnelConfig(tc.Input)
			if actual != tc.Expected {
				t.Fatal("Failed")
			}
			if err == nil && tc.shouldError {
				t.Fatal("Failed")
			}
		})
	}
}
