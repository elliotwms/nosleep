package allfiles

import "time"

// With -all-files, sleeps outside test files are reported too.
func productionSleep() {
	time.Sleep(time.Second) // want "time.Sleep detected"
}

func allowedInProduction() {
	time.Sleep(time.Second) //nosleep:allow polling a device which has no interrupt line
}

// unusedInProduction is reported under -all-files, because this file is
// analysed and so the directive really does govern nothing.
func unusedInProduction() {
	// want +1 `governs no time.Sleep call`
	//nosleep:allow the sleep this justified has been removed
	_ = 1
}
