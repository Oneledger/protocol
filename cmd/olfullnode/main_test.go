// This test file is for global setup and teardown for tests to place their artifacts

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var dir string

func TestMain(m *testing.M) {
	setup := func() {
		// Create a temp file
		var err error
		dir, err = ioutil.TempDir("", "olfullnode")
		if err != nil {
			panic("Failed to create temporary directory: " + dir)
		}
	}

	teardown := func() {
		_ = os.RemoveAll(dir)
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
