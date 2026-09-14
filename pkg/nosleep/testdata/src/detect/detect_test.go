package detect

import (
	"time"
	. "time"
	clock "time"
)

func direct() {
	time.Sleep(time.Second) // want "time.Sleep detected"
}

// aliased is missed by a matcher which compares the package identifier to the
// string "time".
func aliased() {
	clock.Sleep(clock.Second) // want "time.Sleep detected"
}

// dotImported has no package qualifier at all.
func dotImported() {
	Sleep(Second) // want "time.Sleep detected"
}

type sleeper struct{}

func (sleeper) Sleep(d time.Duration) {}

// shadowed is a false positive for a matcher which compares identifiers: the
// receiver is named time but has nothing to do with the time package.
func shadowed() {
	time := sleeper{}
	time.Sleep(0)
}

// indirect is a known gap. Resolving a call through a function value would
// require tracking the assignment, which is not worth the complexity.
var indirect = time.Sleep

func viaVariable() {
	indirect(time.Second)
}
