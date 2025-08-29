package helpers

import (
	"os"
	"testing"
)

var helpers = Helpers{}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
