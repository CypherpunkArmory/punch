// +build all unit

package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/kami-zh/go-capturer"
)

func TestCheckSubdomain(t *testing.T) {
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
			actual := checkSubdomain(tc.Input)
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
