package pkg

// A callback that converts a PendingDevice into a SmartDevice.
type RegisterDevice struct {
	device       PendingDevice
	onRegistered OnRegistered
}

func (cb RegisterDevice) Run(c *Client, m *AppResponse) {
	// Set out if none provided
	if cb.device.out == nil {
		cb.device.out = &SmartDevice{}
	}

	if m.EntityInfo != nil {
		smartDevice := NewSmartDevice(cb.device.GetId(), cb.device.GetName(), *m.EntityInfo.Type)
		c.AddDevice(smartDevice)

		if cb.device.out != nil {
			// Copy device pointer to out
			*cb.device.out = *smartDevice
		}

		if cb.device.onRegistered != nil {
			// Run the followup callback
			cb.onRegistered(cb.device.out)
		}

		if cb.device.outChan != nil {
			// Send the device back to the channel
			cb.device.outChan <- cb.device.out
		}
	}
}
