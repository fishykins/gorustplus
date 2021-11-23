package pkg

type Callback interface {
	Run(c *Client, m *AppResponse)
}

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

// Callback for registering a device
type RegisterCallback struct {
	inner        func(c *Client, m *AppResponse, d PendingDevice, out *SmartDevice)
	onRegistered OnRegistered
	device       PendingDevice
	out          *SmartDevice
}

func (cb RegisterCallback) Run(c *Client, m *AppResponse) {
	if cb.out == nil {
		cb.out = &SmartDevice{}
	}
	// Run the inital callback to register the device
	cb.inner(c, m, cb.device, cb.out)

	// Run the followup callback if one was provided
	if cb.onRegistered != nil {
		cb.onRegistered(cb.out)
	}
}
