package calleeorder

func run() {
	first()
	second()
}

func second() {} // want `function "second" is called by "run" before "first" but declared after it \(stepdown rule\)`

func first() {}
