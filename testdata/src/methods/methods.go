package methods

type Server struct{}

func (s *Server) handle() {} // want `function "handle" is called by "run" but declared before it \(stepdown rule\)`

func (s *Server) run() {
	s.handle()
}
