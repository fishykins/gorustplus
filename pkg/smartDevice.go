package pkg

import (
	"fmt"
)

// A device with a known type
type SmartDevice struct {
	id         uint32
	name       string
	entityType AppEntityType
	onUpdate   map[uint32]func(d *SmartDevice, b *AppEntityPayload)
	uSeq       uint32
}

func NewSmartDevice(id uint32, name string, entityType AppEntityType) *SmartDevice {
	return &SmartDevice{
		id:         id,
		name:       name,
		entityType: entityType,
		uSeq:       0,
	}
}

func (d *SmartDevice) GetId() uint32 {
	return d.id
}

func (d *SmartDevice) GetName() string {
	return d.name
}

func (d *SmartDevice) GetEntityType() AppEntityType {
	return d.entityType
}

// Registers an update event for this device.
func (d *SmartDevice) AddUpdateEvent(f func(d *SmartDevice, b *AppEntityPayload)) uint32 {
	d.uSeq++
	d.onUpdate[d.uSeq] = f
	return d.uSeq
}

// Removes an update event from the device
func (d *SmartDevice) RemoveUpdateEvent(i uint32) error {
	if _, ok := d.onUpdate[i]; ok {
		delete(d.onUpdate, i)
		return nil
	}
	return fmt.Errorf("event id %d is not registered", i)
}

// Called by broadcast handler
func (d *SmartDevice) Update(b *AppEntityPayload) {
	for _, f := range d.onUpdate {
		f(d, b)
	}
}
