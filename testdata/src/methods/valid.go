package methods

type Client struct{}

func (c *Client) send() {
	c.connect()
}

func (c *Client) connect() {}
