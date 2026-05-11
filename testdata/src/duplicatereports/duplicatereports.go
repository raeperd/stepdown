package duplicatereports

func a() {} // want `function "a" is called by "run" but declared before it \(stepdown rule\)`

func run() {
	b()
	a()
}

func b() {}
