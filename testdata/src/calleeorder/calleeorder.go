package calleeorder

func run() {
	first()
	second()
}

func second() {} // want `function "second" is called by "run" after "first" but declared before it \(stepdown rule\)`

func first() {}
