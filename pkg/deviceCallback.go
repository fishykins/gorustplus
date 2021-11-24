package pkg

import (
	"fmt"
	"os"
)

type DeviceCallbackFunc func(c *Client, m *AppResponse, d Device)

// Standard device callback.
type DeviceCallback struct {
	device   *Device
	callback DeviceCallbackFunc
}

func NewDeviceCallback(device *Device, callback DeviceCallbackFunc) (*DeviceCallback, error) {
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}
	return &DeviceCallback{callback: callback, device: device}, nil
}

func (dcb DeviceCallback) Run(c *Client, m *AppResponse) {
	if dcb.device == nil {
		fmt.Println("DeviceCallback has nil device: This should not happen!")
		os.Exit(2)
		return
	}
	if dcb.callback != nil {
		dcb.callback(c, m, *dcb.device)
	}

	// Update cached values.
	if m.EntityInfo != nil {
		dcb.device.SetData(m.EntityInfo.Payload)
	}
}
