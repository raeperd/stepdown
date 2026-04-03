package basic

func bottomCallee() {} // want `function "bottomCallee" is called by "topCaller" but declared before it \(stepdown rule\)`

func topCaller() {
	bottomCallee()
	middleHelper()
}

func middleHelper() {}
