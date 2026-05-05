package qualified_exclusion

type Server struct{}

// No violation — Server.handle is excluded exactly.
func (s *Server) handle() {}

func (s *Server) run() {
	s.handle()
}

type Client struct{}

func (c *Client) handle() {} // want `function "Client.handle" is called by "Client.run" but declared before it \(stepdown rule\)`

func (c *Client) run() {
	c.handle()
}
