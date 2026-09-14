package nosleep_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/elliotwms/nosleep/pkg/nosleep"
)

// TestAnalyzer covers the default configuration: only test files are checked.
func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), nosleep.NewAnalyzer(), "detect", "directive")
}

// TestAnalyzerAll covers -all-files, under which every file is checked.
func TestAnalyzerAll(t *testing.T) {
	a := nosleep.NewAnalyzer()
	if err := a.Flags.Set("all-files", "true"); err != nil {
		t.Fatalf("set -all-files: %s", err)
	}

	analysistest.Run(t, analysistest.TestData(), a, "allfiles")
}

// TestNewAnalyzerIsolatesFlags guards the reason NewAnalyzer exists: two
// analyzers must not share flag state.
func TestNewAnalyzerIsolatesFlags(t *testing.T) {
	a, b := nosleep.NewAnalyzer(), nosleep.NewAnalyzer()

	if err := a.Flags.Set("all-files", "true"); err != nil {
		t.Fatalf("set -all-files: %s", err)
	}

	if got := b.Flags.Lookup("all-files").Value.String(); got != "false" {
		t.Errorf("second analyzer -all-files = %s, want false", got)
	}
}
