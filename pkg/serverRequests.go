package pkg

func (c *Client) GetServerInfo(callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetInfo = &AppEmpty{}
	return c.sendRequest(request, callback)
}

func (c *Client) GetTime(callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetTime = &AppEmpty{}
	return c.sendRequest(request, callback)
}

func (c *Client) GetMap(callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetMap = &AppEmpty{}
	return c.sendRequest(request, callback)
}

func (c *Client) GetMapMarkers(callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetMapMarkers = &AppEmpty{}
	return c.sendRequest(request, callback)
}

func (c *Client) GetCameraFrame(name string, callback func(*AppResponse)) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.GetCameraFrame = &AppCameraFrameRequest{
		Identifier: &name,
	}
	return c.sendRequest(request, callback)
}
