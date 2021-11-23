package pkg

type DeviceCallbackFunc func(c *Client, m *AppResponse, d SmartDevice)
type OnRegistered func(d *SmartDevice)

// Standard device callback.
type DeviceCallback struct {
	inner  DeviceCallbackFunc
	device *SmartDevice
}

func (cb DeviceCallback) Run(c *Client, m *AppResponse) {
	cb.inner(c, m, *cb.device)
}
