package nosleep_test

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"

	root "github.com/elliotwms/nosleep"
)

// TestRegistered guards the init: golangci-lint finds the plugin by name in
// the registry, so a missing or misspelt registration would build a custom
// binary in which the linter cannot be enabled.
func TestRegistered(t *testing.T) {
	if _, err := register.GetPlugin(root.Name); err != nil {
		t.Fatalf("GetPlugin(%q): %s", root.Name, err)
	}
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		conf         any
		wantErr      bool
		wantAllFiles string
	}{
		"no settings block":    {nil, false, "false"},
		"empty settings":       {map[string]any{}, false, "false"},
		"all-files true":       {map[string]any{"all-files": true}, false, "true"},
		"all-files false":      {map[string]any{"all-files": false}, false, "false"},
		"unknown setting":      {map[string]any{"all_files": true}, true, ""},
		"wrong type":           {map[string]any{"all-files": "yes"}, true, ""},
		"not a settings block": {"all-files", true, ""},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := root.New(test.conf)
			if (err != nil) != test.wantErr {
				t.Fatalf("New(%v) error = %v, wantErr %t", test.conf, err, test.wantErr)
			}

			if test.wantErr {
				return
			}

			if got := p.GetLoadMode(); got != register.LoadModeTypesInfo {
				t.Errorf("GetLoadMode() = %q, want %q", got, register.LoadModeTypesInfo)
			}

			analyzers, err := p.BuildAnalyzers()
			if err != nil {
				t.Fatalf("BuildAnalyzers(): %s", err)
			}

			if len(analyzers) != 1 {
				t.Fatalf("BuildAnalyzers() returned %d analyzers, want 1", len(analyzers))
			}

			if got := analyzers[0].Flags.Lookup("all-files").Value.String(); got != test.wantAllFiles {
				t.Errorf("all-files = %s, want %s", got, test.wantAllFiles)
			}
		})
	}
}
