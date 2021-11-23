package websocket

import (
	"errors"
	"fmt"
	"os"

	"github.com/fishykins/gorust/pkg"
	"github.com/fishykins/gorust/pkg/devices"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	connectionData *ConnectionData
	connection     *websocket.Conn
	seq            uint32
	devices        []*devices.SmartDevice
	callbacks      map[uint32]Callback
}

func NewClient(connectionData *ConnectionData) *Client {
	return &Client{
		connectionData: connectionData,
		seq:            0,
		devices:        make([]*devices.SmartDevice, 0),
		callbacks:      make(map[uint32]Callback),
	}
}

func (c *Client) Connect() error {
	var err error
	c.connection, _, err = websocket.DefaultDialer.Dial(c.connectionData.URL(), nil)
	if err != nil {
		return err
	}
	if len(c.devices) > 0 {
		for _, device := range c.devices {
			c.initDevice(*device)
		}
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

// Adds a device to be registered post connection
func (c *Client) RegisterDevice(device devices.SmartDevice) error {
	c.devices = append(c.devices, &device)
	if c.connection != nil {
		return c.initDevice(device)
	}
	return nil
}

// Force adds a device, bypassing server authentication.
//This is useful for when we already know what kind of device it is and its current state.
func (c *Client) ForceAddDevice(device devices.SmartDevice) error {
	c.devices = append(c.devices, &device)
	return nil
}

func (c *Client) RemoveDevice(d devices.SmartDevice) error {
	for i := 0; i < len(c.devices); i++ {
		device := *c.devices[i]
		if device.GetId() == d.GetId() {
			c.devices = append(c.devices[:i], c.devices[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("device not found: %d", d.GetId())
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
		fmt.Println("Write (with cb):", request)
	} else {
		fmt.Println("Write:", request)
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
		fmt.Println("Read:", &appMessage)
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

// =====================================================================================================================
// ============================================== Private Functions ====================================================
// =====================================================================================================================

func (c *Client) getSeq() uint32 {
	s := c.seq
	c.seq++
	return s
}

// Builds a base request.
func (c *Client) newRequest() (*pkg.AppRequest, error) {
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

// Send a request for device info so we can spesify the device type
func (c *Client) initDevice(device devices.SmartDevice) error {
	id := device.GetId()
	name := device.GetName()
	fmt.Printf("%s (%d) init...\n", name, id)

	req, err := c.newRequest()
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
				fmt.Printf("Device %s is type %s", name, entType.String())
				err := c.RemoveDevice(d)
				if err != nil {
					fmt.Println("Error removing device:", err)
					os.Exit(4)
				}
				switch entType {
				case pkg.AppEntityType_Alarm.Enum():
					{
						d := devices.NewSmartAlarm(id, name)
						c.ForceAddDevice(d)
					}
				case pkg.AppEntityType_Switch.Enum():
					{
						d := devices.NewSmartSwitch(id, name)
						c.ForceAddDevice(d)
					}

				case pkg.AppEntityType_StorageMonitor.Enum().Enum():
					{
						d := devices.NewSmartBox(id, name)
						c.ForceAddDevice(d)
					}
				}
			}
		},
	}
	return c.Write(req, cb)
}
