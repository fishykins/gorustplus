package devices

type SmartDevice interface {
	GetId() uint32
	GetName() string
}
