package devices

type SmartAlarm struct {
	name string
	id   uint32
}

func NewSmartAlarm(id uint32, name string) SmartSwitch {
	return SmartSwitch{
		id:   id,
		name: name,
	}
}

func (s SmartAlarm) GetId() uint32 {
	return s.id
}

func (s SmartAlarm) GetName() string {
	return s.name
}
