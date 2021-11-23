package pkg

import "errors"

func (c *Client) SetDeviceState(device *Device, state bool) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	if device.GetType() != AppEntityType_Switch {
		return errors.New("Device is not a switch")
	}
	request.EntityId = &device.Id
	request.SetEntityValue = &AppSetEntityValue{
		Value: &state,
	}
	return c.sendRequest(request, nil)
}

func (c *Client) GetDeviceState(device *Device) error {
	request, err := c.newRequest()
	if err != nil {
		return err
	}
	request.EntityId = &device.Id
	request.GetEntityInfo = &AppEmpty{}
	err = c.Write(request)

	if err != nil {
		return err
	}
	c.deviceCallbacks[*request.Seq] = device
	return nil
}
