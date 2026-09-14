package allfiles

import "time"

// With -all, sleeps outside test files are reported too.
func productionSleep() {
	time.Sleep(time.Second) // want "time.Sleep detected"
}

func allowedInProduction() {
	time.Sleep(time.Second) //nosleep:allow polling a device which has no interrupt line
}
