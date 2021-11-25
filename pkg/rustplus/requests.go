package rustplus

import "errors"

// A set of helper functions for building requests. Execute these requests using 'Client.Write'.

// Builds a base request.
func (c *Client) NewRequest() (*AppRequest, error) {
	if c.connection == nil {
		return nil, errors.New("connection is nil")
	}
	if len(c.connectionData.Tokens) == 0 {
		return nil, errors.New("no tokens")
	}
	token := c.connectionData.Tokens[0]

	seq := c.GetSeq()
	request := AppRequest{
		Seq:         &seq,
		PlayerId:    &token.SteamId,
		PlayerToken: &token.Token,
	}
	return &request, nil
}

func (c *Client) NewDeviceGetRequest(device *Device) (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	id := device.GetId()
	req.EntityId = &id
	req.GetEntityInfo = &AppEmpty{}
	return req, nil
}

func (c *Client) NewDeviceSetRequest(device *Device, state bool) (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	id := device.GetId()
	req.EntityId = &id
	req.SetEntityValue = &AppSetEntityValue{Value: &state}
	return req, nil
}

func (c *Client) NewInfoRequest() (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetInfo = &AppEmpty{}
	return req, nil
}

func (c *Client) NewTimeRequest() (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetTime = &AppEmpty{}
	return req, nil
}

func (c *Client) NewMapRequest() (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetMap = &AppEmpty{}
	return req, nil
}

func (c *Client) NewTeamRequest() (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetTeamInfo = &AppEmpty{}
	return req, nil
}

func (c *Client) NewChatReadRequest() (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetTeamChat = &AppEmpty{}
	return req, nil
}

func (c *Client) NewChatWriteRequest(message string) (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SendTeamMessage = &AppSendMessage{Message: &message}
	return req, nil
}

func (c *Client) NewMarkersRequest() (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetMapMarkers = &AppEmpty{}
	return req, nil
}

func (c *Client) NewCameraRequest(id string, frame uint32) (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.GetCameraFrame = &AppCameraFrameRequest{Identifier: &id, Frame: &frame}
	return req, nil
}

func (c *Client) NewPromoteRequest(id uint64) (*AppRequest, error) {
	req, err := c.NewRequest()
	if err != nil {
		return nil, err
	}
	req.PromoteToLeader = &AppPromoteToLeader{SteamId: &id}
	return req, nil
}
