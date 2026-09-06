package genericmethods

type Box[T any] struct{}

func (b Box[T]) handle() {} // want `function "Box.handle" is called by "Box.run" but declared before it \(stepdown rule\)`

func (b Box[T]) run() {
	b.handle()
}
