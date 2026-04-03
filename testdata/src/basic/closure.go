package basic

func closureHelper() {}

func register() {
	callback := func() {
		closureHelper()
	}
	_ = callback
}
