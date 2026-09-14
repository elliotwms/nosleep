package allfiles

import "time"

func testSleep() {
	time.Sleep(time.Second) // want "time.Sleep detected"
}
