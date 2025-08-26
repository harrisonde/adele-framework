package badgerdriver

import (
	"log"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v3"
)

var testBadgerCache BadgerCache

func TestMain(m *testing.M) {

	if err := os.RemoveAll("./testdata/tmp/badger"); err != nil {
		log.Fatalf("Failed to remove test directory: %v", err)
	}
	if err := os.MkdirAll("./testdata/tmp/badger", 0755); err != nil {
		log.Fatalf("Failed to create test directories: %v", err)
	}

	db, _ := badger.Open(badger.DefaultOptions("./testdata/tmp/badger"))
	testBadgerCache.Conn = db

	os.Exit(m.Run())
}
