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
		{"Valid", "http:80", tunnelConf{"80", "http"}, false},
		{"Invalid(gibberish)", "kajsbdf&*(&", tunnelConf{}, true},
		{"Invalid(wrong order)", "80:http", tunnelConf{"", ""}, true},
		{"Invalid(Extra values)", "http:80:https", tunnelConf{}, true},
		{"Invalid(wrong forward type)", "chewie:80", tunnelConf{"", ""}, true},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := getTunnelConfig(tc.Input)
			if actual != tc.Expected {
				t.Fatalf("Got %s but Expected %s", actual.String(), tc.Expected.String())
			}
			if err == nil && tc.shouldError {
				t.Fatalf("Got no error but expected %s", err)
			}
		})
	}
}
