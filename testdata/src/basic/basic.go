package basic

func callee() {} // want `function "callee" is called by "caller" but declared before it \(stepdown rule\)`

func caller() {
	callee()
}
