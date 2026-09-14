package directive

import "time"

func sameLine() {
	time.Sleep(time.Second) //nosleep:allow the upstream API has no synchronous variant
}

func lineAbove() {
	//nosleep:allow the upstream API has no synchronous variant
	time.Sleep(time.Second)
}

// spaceAfterSlashes is tolerated even though the canonical form has no space.
func spaceAfterSlashes() {
	time.Sleep(time.Second) // nosleep:allow the upstream API has no synchronous variant
}

// bareDirective does not suppress: the reason is mandatory. It is not also
// reported as unused, since it did govern the call.
func bareDirective() {
	// want +1 `//nosleep:allow requires a reason`
	//nosleep:allow
	time.Sleep(time.Second) // want "time.Sleep detected"
}

// notADirective must not be mistaken for one, and so governs nothing.
func notADirective() {
	time.Sleep(time.Second) //nosleep:allowlist foo // want "time.Sleep detected"
}

// twoLinesAbove is too far away to govern the call, which makes it unused.
func twoLinesAbove() {
	// want +1 `governs no time.Sleep call`
	//nosleep:allow this is not close enough to the call

	time.Sleep(time.Second) // want "time.Sleep detected"
}

// leftBehind is the case unused detection exists for: the sleep is gone but
// the excuse for it was not.
func leftBehind() {
	// want +1 `governs no time.Sleep call`
	//nosleep:allow the flaky fixture was replaced with a readiness check
	_ = 1
}

// bareAndUnused is reported as unused rather than as missing a reason: the
// fix is to delete it, not to write a justification for nothing.
func bareAndUnused() {
	// want +1 `governs no time.Sleep call`
	//nosleep:allow
	_ = 1
}

// twoDirectives checks precedence: the trailing directive governs the call, so
// the one above it governs nothing.
func twoDirectives() {
	// want +1 `governs no time.Sleep call`
	//nosleep:allow this one is shadowed by the trailing directive
	time.Sleep(time.Second) //nosleep:allow this one wins
}
