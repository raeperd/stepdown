package interfacecall

type Worker interface {
	work()
}

func run(w Worker) {
	w.work()
}
