// Package nosleep provides a go/analysis pass which reports calls to
// time.Sleep.
package nosleep

import (
	"flag"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const doc = `nosleep: report calls to time.Sleep

Sleeping to wait for something to happen makes a test slow when the sleep is
too long and flaky when it is too short. nosleep reports every call to
time.Sleep so that each one has to be justified rather than reached for by
habit.

By default only test files are checked, since sleeping in production code is
often legitimate. Pass -all-files to check every file.

A call may be allowed with a directive comment giving the reason it is
necessary, placed on any line of the call or on the line above it:

	time.Sleep(time.Second) //nosleep:allow the API has no synchronous variant

The reason is mandatory: a bare //nosleep:allow is reported, and does not
suppress the diagnostic.

A directive which governs no call is reported too, so that a justification
left behind after its sleep was removed does not sit there looking
load-bearing.`

// URL is the documentation link attached to reported diagnostics.
const URL = "https://github.com/elliotwms/nosleep"

// Analyzer is the nosleep analysis pass with default settings.
//
// It carries package-level flag state, so tests which need to vary the
// configuration should call NewAnalyzer instead of mutating this value.
var Analyzer = NewAnalyzer()

// Settings configures the analyzer. The zero value is the default behaviour.
//
// The JSON tags name the keys accepted under the linter's settings block in a
// golangci-lint configuration, and match the command-line flags of the
// standalone binary.
type Settings struct {
	// AllFiles checks every file rather than only _test.go files.
	AllFiles bool `json:"all-files"`
}

// NewAnalyzer returns a new nosleep analyzer with default settings and its own
// flag state.
func NewAnalyzer() *analysis.Analyzer {
	return NewAnalyzerWithSettings(Settings{})
}

// NewAnalyzerWithSettings returns a new nosleep analyzer configured by s. The
// analyzer's flags are bound to a private copy of s, so flags given on the
// command line override it and the caller's value is left alone.
func NewAnalyzerWithSettings(s Settings) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name:     "nosleep",
		Doc:      doc,
		URL:      URL,
		Run:      s.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	a.Flags.Init("nosleep", flag.ExitOnError)
	a.Flags.BoolVar(&s.AllFiles, "all-files", s.AllFiles, "check all files, not just _test.go files")

	return a
}

func (s *Settings) run(pass *analysis.Pass) (any, error) {
	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		// Unreachable via the analysis framework, which guarantees Requires
		// has run, but a nil inspector would panic below.
		return nil, nil
	}

	ds := newDirectives(pass.Fset, pass.Files)

	analysed := func(filename string) bool {
		return s.AllFiles || isTestFile(filename)
	}

	in.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		call := node.(*ast.CallExpr) // always a CallExpr thanks to the node filter

		if !isTimeSleep(pass.TypesInfo, call) {
			return
		}

		pos := pass.Fset.Position(call.Pos())
		if !analysed(pos.Filename) {
			return
		}

		end := pass.Fset.Position(call.End())
		if d, ok := ds.lookup(pos.Filename, pos.Line, end.Line); ok {
			if d.reason != "" {
				return
			}

			// Fail closed: a directive with no reason suppresses nothing,
			// otherwise the requirement to justify a sleep could be dodged by
			// writing the directive and nothing else.
			pass.Reportf(d.pos, "//%s requires a reason: //%s <reason>", DirectivePrefix, DirectivePrefix)
		}

		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			End:     call.End(),
			Message: "time.Sleep detected: justify it with //" + DirectivePrefix + " <reason> or wait on a condition instead",
			URL:     URL,
		})
	})

	ds.reportUnused(pass, analysed)

	return nil, nil
}

// isTimeSleep reports whether call resolves to time.Sleep.
//
// The callee is resolved through the type system rather than matched on the
// source text, so that aliased and dot imports are caught and a method named
// Sleep on a local variable which happens to be named time is not.
func isTimeSleep(info *types.Info, call *ast.CallExpr) bool {
	fn, ok := typeutil.Callee(info, call).(*types.Func)
	if !ok || fn.Pkg() == nil {
		return false
	}

	return fn.Pkg().Path() == "time" && fn.Name() == "Sleep"
}

// isTestFile reports whether filename is a Go test file.
func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}
