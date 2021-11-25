package rustplus

import (
	"fmt"
	"os"
)

type Callback interface {
	Call(m *AppResponse)
}

//====================================================================================
//============================== Basic Callback ======================================
//====================================================================================

type PrimitiveCallback struct {
	inner func(m *AppResponse)
}

func (cb *PrimitiveCallback) Call(m *AppResponse) {
	cb.inner(m)
}

func NewPrimitiveCb(inner func(m *AppResponse)) *PrimitiveCallback {
	return &PrimitiveCallback{inner}
}

//====================================================================================
//============================== Device Callback =====================================
//====================================================================================

type DeviceCallbackFunc func(m *AppResponse, d Device)

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

func (dcb DeviceCallback) Call(m *AppResponse) {
	if dcb.device == nil {
		fmt.Println("DeviceCallback has nil device: This should not happen!")
		os.Exit(2)
		return
	}
	if dcb.callback != nil {
		dcb.callback(m, *dcb.device)
	}

	// Update cached values.
	if m.EntityInfo != nil {
		dcb.device.SetData(m.EntityInfo.Payload)
	}
}

//====================================================================================
//============================ Server Info Callback ==================================
//====================================================================================

type ServerCallback struct {
	inner func(info *AppInfo)
}

func (cb *ServerCallback) Call(m *AppResponse) {
	if m.Info != nil {
		cb.inner(m.Info)
	}
}

func NewServerCb(inner func(info *AppInfo)) *ServerCallback {
	return &ServerCallback{inner}
}

//====================================================================================
//================================ Map Callback ======================================
//====================================================================================

type MapCallback struct {
	inner func(data *AppMap)
}

func (cb *MapCallback) Call(m *AppResponse) {
	if cb.inner != nil && m.Map != nil {
		cb.inner(m.Map)
	}
}

func NewMapCb(inner func(data *AppMap)) *MapCallback {
	return &MapCallback{inner}
}
