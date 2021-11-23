package devices

type SmartBox struct {
	name string
	id   uint32
}

func NewSmartBox(id uint32, name string) SmartBox {
	return SmartBox{
		id:   id,
		name: name,
	}
}
func (s SmartBox) GetId() uint32 {
	return s.id
}

func (s SmartBox) GetName() string {
	return s.name
}
