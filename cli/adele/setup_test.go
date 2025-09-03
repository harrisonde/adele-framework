package main

import (
	"os"
	"testing"
)

var testDir = "testdata"

func TestMain(m *testing.M) {

	os.Mkdir("./testdata", 0755)

	// Clean up
	//os.RemoveAll(testDir)
	defer os.RemoveAll("./testdata")

	os.Exit(m.Run())

}
