package pkg

// A device with a known type
type SmartDevice struct {
	id         uint32
	name       string
	entityType AppEntityType
}

func NewSmartDevice(id uint32, name string, entityType AppEntityType) *SmartDevice {
	return &SmartDevice{
		id:         id,
		name:       name,
		entityType: entityType,
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
