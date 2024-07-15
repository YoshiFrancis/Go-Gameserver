package wsserver

func (c *Client) handlePrompt(message string) {
	if c.prompt == USERNAME {
		c.username = message
		c.prompt = NONE
	}
}
