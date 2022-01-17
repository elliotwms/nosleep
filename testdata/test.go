package testdata

import "time"

func doesNotSleep() {}

func sleeps() {
	time.Sleep(1 * time.Nanosecond) // want "time.Sleep detected"
}
