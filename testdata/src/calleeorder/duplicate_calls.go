package calleeorder

func process() {
	validate()
	transform()
	validate() // second call — should not affect ordering
}

func validate() {}
func transform() {}
