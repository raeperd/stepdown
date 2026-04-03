package crossstruct

type Logger struct{}

func (l *Logger) log() {} // want `function "log" is called by "process" but declared before it \(stepdown rule\)`

func process(l *Logger) {
	l.log()
}
