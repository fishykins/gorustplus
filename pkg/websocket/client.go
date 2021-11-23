package websocket

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/fishykins/gorust/pkg"
	"github.com/fishykins/gorust/pkg/devices"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type pendingDevice struct {
	id   uint32
	name string
	out  *devices.SmartDevice
}

func (d *pendingDevice) GetName() string {
	return d.name
}

func (d *pendingDevice) GetId() uint32 {
	return d.id
}

type Client struct {
	connectionData *ConnectionData
	connection     *websocket.Conn
	seq            uint32
	pendingDevices []*pendingDevice
	devices        map[uint32]*devices.SmartDevice
	callbacks      map[uint32]Callback
}

func NewClient(connectionData *ConnectionData) *Client {
	return &Client{
		connectionData: connectionData,
		seq:            0,
		devices:        make(map[uint32]*devices.SmartDevice),
		callbacks:      make(map[uint32]Callback),
	}
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
		c.pendingDevices = []*pendingDevice{}
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
func (c *Client) RegisterDevice(id uint32, name string, out devices.SmartDevice) error {
	device := &pendingDevice{
		id:   id,
		name: name,
		out:  &out,
	}
	if c.connection != nil {
		return c.initDevice(device)
	} else {
		c.pendingDevices = append(c.pendingDevices, device)
	}
	return nil
}

// Force adds a device, bypassing server authentication.
//This is useful for when we already know what kind of device it is and want to handle its initialization oursleves.
func (c *Client) AddDevice(device devices.SmartDevice) error {
	fmt.Printf("Added device \"%s\" of type %s\n", device.GetName(), reflect.TypeOf(device))
	c.devices[device.GetId()] = &device
	return nil
}

// Removes a device from the client.
func (c *Client) RemoveDevice(d devices.SmartDevice) error {
	if _, ok := c.devices[d.GetId()]; ok {
		delete(c.devices, d.GetId())
		return nil
	}
	return fmt.Errorf("device not found: %d", d.GetId())
}

func (c *Client) TryGetDevice(id uint32) (*devices.SmartDevice, error) {
	if _, ok := c.devices[id]; ok {
		return c.devices[id], nil
	}
	return nil, fmt.Errorf("device not found: %d", id)
}

// Writes the request to websocket and prep callback if provided.
func (c *Client) Write(request *pkg.AppRequest, callback Callback) error {
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

func (c *Client) Read() error {
	_, message, err := c.connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("connection error: %s", err)
	}

	if message != nil {
		appMessage := pkg.AppMessage{}
		err = proto.Unmarshal(message, &appMessage)
		//fmt.Println("Read:", &appMessage)
		if err != nil {
			return fmt.Errorf("unmarshalling error: %s", err)
		}

		// If response, execute callback
		if appMessage.Response != nil {
			if c.callbacks[*appMessage.Response.Seq] != nil {
				cb := c.callbacks[*appMessage.Response.Seq]
				cb.Run(c, appMessage.Response)
				delete(c.callbacks, *appMessage.Response.Seq)
			}
			return nil
		}
		// If broadcast, handle any updates
		if appMessage.Broadcast != nil {
			return nil
		}
	}
	return nil
}

// Builds a base request.
func (c *Client) NewRequest() (*pkg.AppRequest, error) {
	if c.connection == nil {
		return nil, errors.New("connection is nil")
	}
	if len(c.connectionData.Tokens) == 0 {
		return nil, errors.New("no tokens")
	}
	token := c.connectionData.Tokens[0]

	seq := c.getSeq()
	request := pkg.AppRequest{
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

// Send a request for device info so we can spesify the device type
func (c *Client) initDevice(device devices.SmartDevice) error {
	id := device.GetId()
	name := device.GetName()
	fmt.Printf("%s (%d) init...\n", name, id)

	req, err := c.NewRequest()
	if err != nil {
		return err
	}

	req.EntityId = &id
	req.GetEntityInfo = &pkg.AppEmpty{}

	cb := DeviceCallback{
		device: &device,
		inner: func(c *Client, m *pkg.AppResponse, d devices.SmartDevice) {
			if m.EntityInfo != nil {
				entType := m.EntityInfo.Type
				id := d.GetId()
				name := d.GetName()
				if err != nil {
					fmt.Println("Error removing device:", err)
					os.Exit(4)
				}
				switch *entType {
				case pkg.AppEntityType_Alarm:
					d := devices.NewSmartAlarm(id, name)
					c.AddDevice(d)
				case pkg.AppEntityType_StorageMonitor:
					d := devices.NewSmartBox(id, name)
					c.AddDevice(d)
				case pkg.AppEntityType_Switch:
					d := devices.NewSmartSwitch(id, name)
					c.AddDevice(d)
				default:
					fmt.Printf("Unknown device type: %s\n", *entType)
				}
			}
		},
	}
	return c.Write(req, cb)
}
