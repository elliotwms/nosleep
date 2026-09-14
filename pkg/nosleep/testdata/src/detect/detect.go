package detect

import "time"

// productionSleep is not reported: only test files are checked by default.
func productionSleep() {
	time.Sleep(time.Second)
}
