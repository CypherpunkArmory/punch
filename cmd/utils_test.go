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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/kami-zh/go-capturer"
)

func TestcorrectSubdomainRegex(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{"Valid", "testdomain", true},
		{"Invalid", " ooiqhwe&*&", false},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := correctSubdomainRegex(tc.Input)
			if actual != tc.Expected {
				t.Fatal("Failed")
			}
		})
	}
}
func TestCheckPort(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{"TooLow", "0", false},
		{"Valid", "12374", true},
		{"TooHigh", "65536", false},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := checkPort(tc.Input)
			if actual != tc.Expected {
				t.Fatal("Failed")
			}
		})
	}
}
func TestFixFilePath(t *testing.T) {
	cases := []struct {
		Name       string
		Input      string
		Expected   string
		ShouldFail bool
	}{
		{"No change needed", "/Users/test", "/Users/test", false},
		{"Append home directory", "~/.ssh", "~/.ssh", true},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := fixFilePath(tc.Input)
			if actual != tc.Expected {
				if tc.ShouldFail {
				} else {
					t.Fatal("Failed")
				}
			}
		})
	}
}
func TestGetPublicKey(t *testing.T) {
	cases := []struct {
		Name       string
		Input      string
		ShouldFail bool
	}{
		{"Valid Pub key", filepath.Join("test-files", "test.pub"), false},
		{"Invalid Pub key", filepath.Join("test-files", "test2.pub"), true},
		{"Incorrect path", filepath.Join("test-files", "test3.pub"), true},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := getPublicKey(tc.Input)
			if actual == "" && !tc.ShouldFail {
				t.Fatal("Failed")
			}
			if err != nil {
				if tc.ShouldFail {
				} else {
					t.Fatal("Failed")
				}
			}
		})
	}
}
func TestPrintError(t *testing.T) {
	nonNilErrorMessage := capturer.CaptureStderr(func() {
		reportError("Test", false)
	})
	nilErrorMessage := capturer.CaptureStderr(func() {
		reportError("", false)
	})
	fmt.Println(nonNilErrorMessage)
	fmt.Println(nilErrorMessage)

	// Output:
	// Test
	// Unexpected error occured
}
