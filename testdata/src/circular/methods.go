package circular

type Worker struct{}

func (w *Worker) start() {
	w.stop()
}

func (w *Worker) stop() {
	w.start()
}
