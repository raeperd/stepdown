package short_method_exclusion

type Server struct{}

func (s *Server) handle() {}

func (s *Server) run() {
	s.handle()
}

type Client struct{}

func (c *Client) handle() {}

func (c *Client) run() {
	c.handle()
}
