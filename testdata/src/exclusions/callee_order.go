package exclusions

// No callee-order violation — init is excluded, so it should not
// participate in ordering checks.
func init() {
	helperA()
	helperB()
}

func helperA() {}
func helperB() {}
