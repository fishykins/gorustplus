package devices

type SmartSwitch struct {
	name string
	id   uint32
}

func NewSmartSwitch(id uint32, name string) SmartSwitch {
	return SmartSwitch{
		id:   id,
		name: name,
	}
}

func (s SmartSwitch) GetId() uint32 {
	return s.id
}

func (s SmartSwitch) GetName() string {
	return s.name
}
