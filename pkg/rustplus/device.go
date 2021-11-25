package rustplus

import (
	"fmt"
)

// A device with a known type
type Device struct {
	id         uint32
	Name       string
	value      *bool
	entityType *AppEntityType
	onInit     DeviceCallbackFunc
	onUpdate   map[uint32]BroadcastEvent
	uSeq       uint32
}

func NewDevice(id uint32, name string) *Device {
	return &Device{
		id:         id,
		Name:       name,
		value:      nil,
		entityType: nil,
		onInit:     nil,
		onUpdate:   make(map[uint32]BroadcastEvent),
	}
}

func (d *Device) GetId() uint32 {
	return d.id
}

func (d *Device) GetType() AppEntityType {
	return *d.entityType
}

func (d *Device) SetType(t AppEntityType) error {
	if d.entityType != nil {
		return fmt.Errorf("device type already set")
	}
	d.entityType = &t
	return nil
}

func (d *Device) SetInit(callback DeviceCallbackFunc) {
	d.onInit = callback
}

func (d *Device) GetValue() bool {
	return *d.value
}

// Registers a bradcast event for this device.
func (d *Device) AddBroadcastEvent(f BroadcastEvent) uint32 {
	d.uSeq++
	d.onUpdate[d.uSeq] = f
	return d.uSeq
}

// Removes an update event from the device
func (d *Device) RemoveUpdateEvent(i uint32) error {
	if _, ok := d.onUpdate[i]; ok {
		delete(d.onUpdate, i)
		return nil
	}
	return fmt.Errorf("event id %d is not registered", i)
}

// Called by broadcast handler before values are set.
func (d *Device) BroadcastEvent(b *AppEntityPayload) {
	if d.entityType == nil {
		for _, f := range d.onUpdate {
			f(d, b)
		}
		return
	}
	// Prevent update spamming from Storagebox, as it sends the same update twice with both true and false.
	if *d.entityType != AppEntityType_StorageMonitor || (*d.entityType == AppEntityType_StorageMonitor && *b.Value) {
		for _, f := range d.onUpdate {
			f(d, b)
		}
	}
}

// Sets the device values using websocket payload. Called after BroadcastEvent.
func (d *Device) SetData(b *AppEntityPayload) {
	d.value = b.Value
	// TODO: Add box support
}

// A quick and simple way to write value to the websocket with no callback.
func (d *Device) WriteValue(value bool) error {
	c := GetClient()
	return c.SetDeviceInfo(d, value, nil)
}

func (d *Device) ReadValue(callback DeviceCallbackFunc) error {
	c := GetClient()
	return c.GetDeviceInfo(d, callback)
}
