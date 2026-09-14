package nosleep

import (
	"go/ast"
	"go/token"
	"strings"
)

// DirectivePrefix is the comment directive which suppresses a nosleep
// diagnostic. It must be followed by a non-empty reason:
//
//	time.Sleep(time.Second) //nosleep:allow waiting for the debouncer to settle
//
// A directive applies to the line it appears on, or to the line immediately
// following it.
const DirectivePrefix = "nosleep:allow"

// directive is a parsed //nosleep:allow comment.
type directive struct {
	pos    token.Pos
	reason string
}

// directives indexes the //nosleep:allow directives in a set of files by
// filename and line number.
type directives map[string]map[int]directive

// newDirectives scans files for //nosleep:allow comments and indexes them by
// position. Where two directives share a line the last one wins; that only
// happens in pathological source, and the choice is arbitrary either way.
func newDirectives(fset *token.FileSet, files []*ast.File) directives {
	ds := directives{}

	for _, f := range files {
		for _, group := range f.Comments {
			for _, c := range group.List {
				reason, ok := parseDirective(c.Text)
				if !ok {
					continue
				}

				pos := fset.Position(c.Slash)
				if ds[pos.Filename] == nil {
					ds[pos.Filename] = map[int]directive{}
				}
				ds[pos.Filename][pos.Line] = directive{pos: c.Slash, reason: reason}
			}
		}
	}

	return ds
}

// lookup returns the directive governing the given line, which is either a
// directive on the line itself or one on the line immediately above it.
func (ds directives) lookup(filename string, line int) (directive, bool) {
	byLine, ok := ds[filename]
	if !ok {
		return directive{}, false
	}

	// A trailing directive on the same line takes precedence over one on the
	// line above, so that the closest directive to the call always wins.
	if d, ok := byLine[line]; ok {
		return d, true
	}

	d, ok := byLine[line-1]

	return d, ok
}

// parseDirective reports whether text is a //nosleep:allow comment and, if so,
// returns the reason given for it. The reason is empty when none was given.
func parseDirective(text string) (reason string, ok bool) {
	// Canonically the directive has no space after the slashes, matching
	// //go:build and //nolint, but a space is tolerated since it reads more
	// naturally and is an easy mistake to make.
	body, ok := strings.CutPrefix(text, "//")
	if !ok {
		return "", false
	}

	rest, ok := strings.CutPrefix(strings.TrimSpace(body), DirectivePrefix)
	if !ok {
		return "", false
	}

	// Require the prefix to be a whole word, so that a comment such as
	// //nosleep:allowlist is not mistaken for a directive.
	if rest != "" && rest[0] != ' ' && rest[0] != '\t' {
		return "", false
	}

	return strings.TrimSpace(rest), true
}
