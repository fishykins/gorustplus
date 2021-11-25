package pkg

// Helper functions to handle both request building and writing, with callback support.
// Use the simple request functions if you require greater finesse in your callbacks!

func (c *Client) GetDeviceInfo(d *Device, callback DeviceCallbackFunc) error {
	// Register the device if we need to
	if _, err := c.TryGetDevice(d.GetId()); err != nil {
		c.AddDevice(d)
	}
	request, err := c.NewDeviceGetRequest(d)
	if err != nil {
		return err
	}
	cb, err := NewDeviceCallback(d, callback)
	if err != nil {
		return err
	}
	c.Write(request, cb)
	return nil
}

func (c *Client) SetDeviceInfo(d *Device, state bool, callback DeviceCallbackFunc) error {
	// Register the device if we need to
	if _, err := c.TryGetDevice(d.GetId()); err != nil {
		c.AddDevice(d)
	}
	request, err := c.NewDeviceSetRequest(d, state)
	if err != nil {
		return err
	}
	cb, err := NewDeviceCallback(d, callback)
	if err != nil {
		return err
	}
	c.Write(request, cb)
	return nil
}

func (c *Client) GetServerInfo(callback func(info *AppInfo)) error {
	request, err := c.NewMapRequest()
	if err != nil {
		return err
	}
	cb := NewServerCb(callback)
	if err != nil {
		return err
	}
	c.Write(request, cb)
	return nil
}

func (c *Client) GetMapInfo(callback func(data *AppMap)) error {
	request, err := c.NewMapRequest()
	if err != nil {
		return err
	}
	cb := NewMapCb(callback)
	if err != nil {
		return err
	}
	c.Write(request, cb)
	return nil
}
