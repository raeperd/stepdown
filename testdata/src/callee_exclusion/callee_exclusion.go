package callee_exclusion

// No violation — "excluded" is excluded as callee
func excluded() {}

func caller() {
	excluded()
}

// Non-excluded pair still produces violation
func other() {} // want `function "other" is called by "another" but declared before it \(stepdown rule\)`

func another() {
	other()
}
