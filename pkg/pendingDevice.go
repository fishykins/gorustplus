package pkg

// A device of unknown type
type PendingDevice struct {
	id           uint32
	name         string
	out          *SmartDevice
	onRegistered OnRegistered
}

func (d *PendingDevice) GetId() uint32 {
	return d.id
}

func (d *PendingDevice) GetName() string {
	return d.name
}

func (d *PendingDevice) GetOut() *SmartDevice {
	return d.out
}

func (d *PendingDevice) GetCb() *OnRegistered {
	return &d.onRegistered
}

func NewPendingDevice(id uint32, name string, out *SmartDevice, onRegistered OnRegistered) *PendingDevice {
	return &PendingDevice{
		id:           id,
		name:         name,
		out:          out,
		onRegistered: onRegistered,
	}
}
