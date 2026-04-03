package exclusions

// No violations — init is excluded as caller
func setup() {}

func init() {
	setup()
}

// Non-excluded functions still produce violations
func callee() {} // want `function "callee" is called by "nonExcluded" but declared before it \(stepdown rule\)`

func nonExcluded() {
	callee()
}
