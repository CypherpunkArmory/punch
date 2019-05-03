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

package tunnel

import (
	"path/filepath"
	"testing"
)

func Test_privateKeyFile(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		shouldFail bool
	}{
		{"Valid priv key", filepath.Join("test-files", "test.pem"), false},
		{"Invalid priv key", filepath.Join("test-files", "test2.pem"), true},
		{"Incorrect path", filepath.Join("test-files", "test3.pem"), true},
		{"Empty file", filepath.Join("test-files", "test4.pem"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := readPrivateKeyFile(tc.path)
			if actual == nil && !tc.shouldFail {
				t.Fatal("Failed")
			}
			if actual != nil && tc.shouldFail {
				t.Fatal("Failed")
			}
		})
	}
}
