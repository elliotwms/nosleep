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

// bareDirective does not suppress: the reason is mandatory.
func bareDirective() {
	// want +1 `//nosleep:allow requires a reason`
	//nosleep:allow
	time.Sleep(time.Second) // want "time.Sleep detected"
}

// notADirective must not be mistaken for one.
func notADirective() {
	time.Sleep(time.Second) //nosleep:allowlist foo // want "time.Sleep detected"
}

// twoLinesAbove is too far away to suppress the call.
func twoLinesAbove() {
	//nosleep:allow this is not close enough to the call

	time.Sleep(time.Second) // want "time.Sleep detected"
}
