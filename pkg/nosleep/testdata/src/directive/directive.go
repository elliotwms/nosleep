package directive

// notAnalysedByDefault holds a directive which governs nothing, but is not
// reported: this file is only examined under -all-files, so the directive is
// not idle by its own doing.
func notAnalysedByDefault() {
	//nosleep:allow nothing in this file is analysed by default
	_ = 1
}
