package pkg

func (c *Client) GetTeamInfo(callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetTeamInfo = &AppEmpty{}
	return c.sendRequest(request, callback)
}

// Not working?
func (c *Client) GetTeamMessages(callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetTeamChat = &AppEmpty{}
	return c.sendRequest(request, callback)
}

func (c *Client) SendTeamMessage(message string, callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.SendTeamMessage = &AppSendMessage{
		Message: &message,
	}
	return c.sendRequest(request, callback)
}
