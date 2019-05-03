// Punch CLI used for interacting with holepunch.io
// Copyright (C) 2018-2019  Orb.House, LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
