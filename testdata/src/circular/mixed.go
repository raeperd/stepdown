package circular

func helper() {} // want `function "helper" is called by "entry" but declared before it \(stepdown rule\)`

func entry() {
	helper()
	mutual()
}

func mutual() {
	entry()
}
