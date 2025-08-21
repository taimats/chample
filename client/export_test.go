package client

func ClientRead(c *Client[struct{}]) error {
	return c.read()
}

func ClientWrite(c *Client[struct{}]) error {
	return c.write()
}

func ClientDone(c *Client[struct{}]) {
	c.done <- struct{}{}
}

func ClientSend(c *Client[struct{}], msg *Message) {
	c.send <- msg
}
