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
			actual, _ := privateKeyFile(tc.path)
			if actual == nil && !tc.shouldFail {
				t.Fatal("Failed")
			}
			if actual != nil && tc.shouldFail {
				t.Fatal("Failed")
			}
		})
	}
}
