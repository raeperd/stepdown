package circular

func x() {
	y()
}

func y() {
	z()
}

func z() {
	x()
}
