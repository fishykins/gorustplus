package pkg

import "errors"

// A helper struct to store device data
type Device struct {
	Id          uint32
	Name        string
	Description string
	OnUpdate    func(d *Device)
	value       bool
	items       []*AppEntityPayload_Item
	capacity    *int32
	entType     AppEntityType
}

func (d *Device) GetState() (bool, error) {
	if d.value {
		return d.value, nil
	}
	return false, errors.New("Device value is not set")
}

func (d *Device) GetItems() ([]*AppEntityPayload_Item, error) {
	if d.items != nil {
		return d.items, nil
	}
	return nil, errors.New("Device items are not set")
}

func (d *Device) GetCapacity() (int32, error) {
	if d.capacity != nil {
		return *d.capacity, nil
	}
	return 0, errors.New("Device capacity is not set")
}

func (d *Device) GetType() AppEntityType {
	return d.entType
}
