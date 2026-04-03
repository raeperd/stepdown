package calleeorder

// No callee-order violation — b is in a cycle with a, so b should not
// set the max-line reference for ordering helper.
func a2() {
	b2()
	helper2()
}

func helper2() {}

func b2() {
	a2()
}
