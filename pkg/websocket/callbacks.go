package websocket

import (
	"github.com/fishykins/gorust/pkg"
	"github.com/fishykins/gorust/pkg/devices"
)

type Callback interface {
	Run(c *Client, m *pkg.AppResponse)
}

type DeviceCallback struct {
	inner  func(c *Client, m *pkg.AppResponse, d devices.SmartDevice)
	device *devices.SmartDevice
}

func (cb DeviceCallback) Run(c *Client, m *pkg.AppResponse) {
	cb.inner(c, m, *cb.device)
}
