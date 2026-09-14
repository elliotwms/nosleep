package nosleep_test

import (
	"flag"
	"testing"

	"github.com/elliotwms/nosleep/pkg/nosleep"
)

// reservedFlags are the flag names that singlechecker and multichecker register
// on the global flag set, including the deprecated vet shims. An analyzer flag
// which collides with one of these panics at startup, before any analysis runs,
// so the failure is invisible to analysistest and only shows up when the binary
// is actually invoked.
//
// Sourced from golang.org/x/tools/go/analysis/internal/analysisflags.
var reservedFlags = []string{
	"V", "all", "c", "diff", "fix", "flags", "json", "source", "tags", "v",
	"bool", "buildtags", "methods", "rangeloops",
	"compositewhitelist", "printfuncs", "shadowstrict",
	"unusedfuncs", "unusedstringmethods",
}

func TestFlagsDoNotCollideWithChecker(t *testing.T) {
	reserved := make(map[string]bool, len(reservedFlags))
	for _, name := range reservedFlags {
		reserved[name] = true
	}

	nosleep.NewAnalyzer().Flags.VisitAll(func(f *flag.Flag) {
		if reserved[f.Name] {
			t.Errorf("flag -%s collides with a name registered by singlechecker; "+
				"the binary will panic on startup", f.Name)
		}
	})
}
