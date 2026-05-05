package methods

type Server struct{}

func (s *Server) handle() {} // want `function "Server.handle" is called by "Server.run" but declared before it \(stepdown rule\)`

func (s *Server) run() {
	s.handle()
}
