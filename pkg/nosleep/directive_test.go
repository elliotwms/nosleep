package nosleep

import "testing"

func TestParseDirective(t *testing.T) {
	tests := map[string]struct {
		text       string
		wantOK     bool
		wantReason string
	}{
		"directive with reason":   {"//nosleep:allow because reasons", true, "because reasons"},
		"space after slashes":     {"// nosleep:allow because reasons", true, "because reasons"},
		"tab before reason":       {"//nosleep:allow\tbecause reasons", true, "because reasons"},
		"trailing whitespace":     {"//nosleep:allow because reasons  ", true, "because reasons"},
		"bare directive":          {"//nosleep:allow", true, ""},
		"bare with trailing tab":  {"//nosleep:allow\t", true, ""},
		"longer word":             {"//nosleep:allowlist foo", false, ""},
		"unrelated comment":       {"// just a comment", false, ""},
		"mentions but not prefix": {"// see nosleep:allow for details", false, ""},
		"block comment":           {"/*nosleep:allow because reasons*/", false, ""},
		"empty":                   {"", false, ""},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reason, ok := parseDirective(test.text)

			if ok != test.wantOK {
				t.Fatalf("parseDirective(%q) ok = %t, want %t", test.text, ok, test.wantOK)
			}

			if reason != test.wantReason {
				t.Errorf("parseDirective(%q) reason = %q, want %q", test.text, reason, test.wantReason)
			}
		})
	}
}
