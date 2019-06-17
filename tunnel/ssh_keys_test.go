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
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_privateKeyFile(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		shouldFail bool
	}{
		{"Valid RSA priv key", filepath.Join("test-files", "test.pem"), false},
		{"Valid OPENSSH priv key", filepath.Join("test-files", "test5.pem"), false},
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
func Test_privatePasswordProtectedKeyFile(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		password   string
		shouldFail bool
	}{
		{"Valid OPENSSH priv key w/ password", filepath.Join("test-files", "test6.pem"), "test", false},
		{"Valid OPENSSH priv key w/ wrong password", filepath.Join("test-files", "test6.pem"), "wrong pass", true},
		{"Valid RSA priv key w/ password", filepath.Join("test-files", "test7.pem"), "test", false},
		{"Valid RSA priv key w/ wrong password", filepath.Join("test-files", "test7.pem"), "wrong pass", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buffer, _ := ioutil.ReadFile(tc.path)
			actual, _ := readPasswordProtectedKey(buffer, []byte(tc.password), tc.path)
			if actual == nil && !tc.shouldFail {
				t.Fatal("Failed")
			}
			if actual != nil && tc.shouldFail {
				t.Fatal("Failed")
			}
		})
	}
}
