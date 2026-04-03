package basic

func c() {} // want `function "c" is called by "b" but declared before it \(stepdown rule\)`

func b() { // want `function "b" is called by "a" but declared before it \(stepdown rule\)`
	c()
}

func a() {
	b()
}
