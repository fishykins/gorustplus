package pkg

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var client = Client{}

type Client struct {
	connectionData *ConnectionData
	connection     *websocket.Conn
	seq            uint32
	pendingDevices []*PendingDevice
	devices        map[uint32]*SmartDevice
	callbacks      map[uint32]Callback
}

// Instansiates a new static client. This will overwrite any existing client so be careful!
func NewClient(connectionData *ConnectionData) *Client {
	client = Client{
		connectionData: connectionData,
		seq:            0,
		devices:        make(map[uint32]*SmartDevice),
		callbacks:      make(map[uint32]Callback),
	}
	return &client
}

// Gets the static client
func GetClient() *Client {
	return &client
}

func (c *Client) Connect() error {
	var err error
	c.connection, _, err = websocket.DefaultDialer.Dial(c.connectionData.URL(), nil)
	if err != nil {
		return err
	}
	if len(c.pendingDevices) > 0 {
		for _, device := range c.pendingDevices {
			c.initDevice(device)
		}
		c.pendingDevices = []*PendingDevice{}
	}
	return nil
}

func (c *Client) Disconnect() error {
	err := c.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

// Adds a device to be registered post connection. This is useful for devices that are of an unknown type,
// or for devices that we want the client to initialize with their current values.
func (c *Client) RegisterDevice(id uint32, name string, out *SmartDevice) error {
	device := NewPendingDevice(id, name, out, nil)
	if c.connection != nil {
		return c.initDevice(device)
	} else {
		c.pendingDevices = append(c.pendingDevices, device)
	}
	return nil
}

// Adds a device to be registered post connection. This is useful for devices that are of an unknown type,
// or for devices that we want the client to initialize with their current values.
func (c *Client) RegisterDeviceWithCallback(id uint32, name string, callback OnRegistered) error {
	device := NewPendingDevice(id, name, nil, callback)
	if c.connection != nil {
		return c.initDevice(device)
	} else {
		c.pendingDevices = append(c.pendingDevices, device)
	}
	return nil
}

// Force adds a device, bypassing server authentication.
//This is useful for when we already know what kind of device it is and want to handle its initialization ourselves.
func (c *Client) AddDevice(device *SmartDevice) error {
	fmt.Printf("Added device \"%s\" of type %s\n", device.GetName(), reflect.TypeOf(device))
	c.devices[device.GetId()] = device
	return nil
}

// Removes a device from the client.
func (c *Client) RemoveDevice(d SmartDevice) error {
	if _, ok := c.devices[d.GetId()]; ok {
		delete(c.devices, d.GetId())
		return nil
	}
	return fmt.Errorf("device not found: %d", d.GetId())
}

func (c *Client) TryGetDevice(id uint32) (*SmartDevice, error) {
	if _, ok := c.devices[id]; ok {
		return c.devices[id], nil
	}
	return nil, fmt.Errorf("device not found: %d", id)
}

// Writes the request to websocket and prep callback if provided.
func (c *Client) Write(request *AppRequest, callback Callback) error {
	data, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	c.connection.WriteMessage(websocket.BinaryMessage, data)
	if callback != nil {
		c.callbacks[*request.Seq] = callback
	}
	return nil
}

// Reads a message from the websocket. This is a blocking call.
func (c *Client) Read() (*AppMessage, error) {
	_, message, err := c.connection.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("connection error: %s", err)
	}

	if message != nil {
		appMessage := AppMessage{}
		err = proto.Unmarshal(message, &appMessage)
		//fmt.Println("Read:", &appMessage)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling error: %s", err)
		}
		return &appMessage, nil
	}
	return nil, nil
}

// Handles the given message, executing any callbacks.
func (c *Client) HandleMessage(message *AppMessage) error {
	if message == nil {
		return errors.New("message is nil")
	}
	// If response, execute callback
	if message.Response != nil {
		return c.handleResponse(message.Response)
	}
	// If broadcast, handle any updates
	if message.Broadcast != nil {
		return c.handleBroadcast(message.Broadcast)
	}
	return nil
}

// Builds a base request.
func (c *Client) NewRequest() (*AppRequest, error) {
	if c.connection == nil {
		return nil, errors.New("connection is nil")
	}
	if len(c.connectionData.Tokens) == 0 {
		return nil, errors.New("no tokens")
	}
	token := c.connectionData.Tokens[0]

	seq := c.getSeq()
	request := AppRequest{
		Seq:         &seq,
		PlayerId:    &token.SteamId,
		PlayerToken: &token.Token,
	}
	return &request, nil
}

// =====================================================================================================================
// ============================================== Private Functions ====================================================
// =====================================================================================================================

func (c *Client) getSeq() uint32 {
	s := c.seq
	c.seq++
	return s
}

func (c *Client) handleResponse(r *AppResponse) error {
	if r == nil {
		return errors.New("response is nil")
	}
	if c.callbacks[*r.Seq] != nil {
		cb := c.callbacks[*r.Seq]
		cb.Run(c, r)
		delete(c.callbacks, *r.Seq)
	}
	return nil
}

func (c *Client) handleBroadcast(b *AppBroadcast) error {
	if b == nil {
		return errors.New("broadcast is nil")
	}
	if b.EntityChanged != nil {
		entityId := *b.EntityChanged.EntityId
		if device, ok := c.devices[entityId]; ok {
			device.Update(b.EntityChanged.Payload)
		} else {
			return fmt.Errorf("device not found: %d", entityId)
		}
	}
	if b.TeamMessage != nil {
		fmt.Println("WARNING: message handling not yet implimented!")
		fmt.Printf("Team message: %s\n", b.TeamMessage.Message)
	}
	if b.TeamChanged != nil {
		player := b.TeamChanged.PlayerId
		teamInfo := b.TeamChanged.TeamInfo
		fmt.Println("WARNING: Team update handling not yet implimented!")
		fmt.Printf("Team update: %v (%v members)\n", player, len(teamInfo.Members))
	}
	return nil
}

// Send a request for device info so we can spesify the device type
func (c *Client) initDevice(device *PendingDevice) error {
	id := device.GetId()
	name := device.GetName()
	fmt.Printf("%s (%d) init...\n", name, id)

	req, err := c.NewRequest()
	if err != nil {
		return err
	}

	req.EntityId = &id
	req.GetEntityInfo = &AppEmpty{}

	cb := RegisterCallback{
		device:       *device,
		out:          device.GetOut(),
		onRegistered: *device.GetCb(),
		inner: func(c *Client, m *AppResponse, d PendingDevice, out *SmartDevice) {
			if m.EntityInfo != nil {
				id := d.GetId()
				name := d.GetName()
				if err != nil {
					fmt.Println("Error removing device:", err)
					os.Exit(4)
				}
				d := NewSmartDevice(id, name, *m.EntityInfo.Type)
				c.AddDevice(d)
				if out != nil {
					*out = *d
				}
			}
		},
	}
	return c.Write(req, cb)
}
