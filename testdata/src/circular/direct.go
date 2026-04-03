package circular

// No violations — a and b call each other
func a() {
	b()
}

func b() {
	a()
}
