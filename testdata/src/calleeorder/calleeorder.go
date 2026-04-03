package calleeorder

func main() {
	first()
	second()
}

func second() {} // want `function "second" is called by "main" before "first" but declared after it \(stepdown rule\)`

func first() {}
