package interfaces

// Calls through an interface must not be reported. The interface type is
// declared in this file, but the method has no body here to move.
type Speaker interface {
	Speak()
}

func viaInterface(s Speaker) {
	s.Speak()
}

// Interfaces declared after their use are not reported either.
func viaLaterInterface(w Writer) {
	w.Write()
}

type Writer interface {
	Write()
}

// Methods with a body in this file are still checked.
type Dog struct{}

func (Dog) bark() {} // want `function "Dog.bark" is called by "Dog.Speak" but declared before it \(stepdown rule\)`

func (Dog) Speak() {
	Dog{}.bark()
}
