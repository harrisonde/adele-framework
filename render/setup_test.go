package render

import (
	"os"
	"testing"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./testdata/views"),
	jet.InDevelopmentMode(),
)

var testRenderer = Render{
	Renderer: "",
	RootPath: "",
	JetViews: views,
}

// Special file, Go will run all tests with the file name of
// setup_test.go and run the function called TestMain().
func TestMain(m *testing.M) {
	os.Exit(m.Run()) // Exit, but before you do, run out tests.
}
