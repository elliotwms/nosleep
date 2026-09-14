// Package nosleep registers the nosleep linter as a golangci-lint module
// plugin.
//
// Build a golangci-lint binary which includes it by listing this module in
// .custom-gcl.yml:
//
//	version: v2.13.2 # the golangci-lint release to build
//	plugins:
//	  - module: github.com/elliotwms/nosleep
//	    version: v1.2.3 # a nosleep release tag, or anything go get accepts
//
// then enable and, optionally, configure it in .golangci.yml. Custom linters
// are not enabled by default, so the enable entry is required:
//
//	linters:
//	  enable:
//	    - nosleep
//	  settings:
//	    custom:
//	      nosleep:
//	        type: module
//	        description: Reports time.Sleep calls that have not been justified.
//	        original-url: github.com/elliotwms/nosleep
//	        settings:
//	          all-files: false
//
// The keys accepted under settings are the JSON tags of [nosleep.Settings].
//
// The analyzer itself lives in package [nosleep] under pkg/nosleep; this
// package only adapts it to the plugin interface, so that importing the
// analyzer does not drag in the plugin registry.
package nosleep

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/elliotwms/nosleep/pkg/nosleep"
)

// Name is the name the linter is registered under, and so the name used to
// enable it in .golangci.yml and in //nolint directives.
const Name = "nosleep"

func init() {
	register.Plugin(Name, New)
}

// New constructs the plugin from the settings block in .golangci.yml, which
// golangci-lint hands over as decoded YAML. Unknown keys are an error rather
// than silently ignored, so that a typo in a setting cannot quietly leave the
// default in place.
func New(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[nosleep.Settings](conf)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", Name, err)
	}

	return &plugin{settings: settings}, nil
}

// plugin adapts the nosleep analyzer to golangci-lint's plugin interface.
type plugin struct {
	settings nosleep.Settings
}

// BuildAnalyzers implements [register.LinterPlugin].
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{nosleep.NewAnalyzerWithSettings(p.settings)}, nil
}

// GetLoadMode implements [register.LinterPlugin]. The analyzer resolves
// callees through the type checker, so it needs type information.
func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
