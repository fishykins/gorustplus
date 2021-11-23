package pkg

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	useProxy        bool
	ip              string
	port            uint64
	connection      *websocket.Conn
	tokens          []PlayerToken
	seq             uint32
	callbacks       map[uint32]func(*AppResponse)
	deviceCallbacks map[uint32]*Device
	devices         map[uint32]*Device
}

func NewClient(ip string, port uint64) Client {
	return Client{
		ip:              ip,
		port:            port,
		callbacks:       make(map[uint32]func(*AppResponse)),
		deviceCallbacks: make(map[uint32]*Device),
		devices:         make(map[uint32]*Device),
	}
}

func (c *Client) Connect() error {
	adr := c.Address()
	connection, _, err := websocket.DefaultDialer.Dial(adr, nil)
	if err != nil {
		return err
	}
	if connection == nil {
		return errors.New("connection is nil")
	}

	c.connection = connection
	c.seq = 0

	for _, device := range c.devices {
		c.initDevice(device)
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

func (c *Client) ClearCache() {
	c.callbacks = make(map[uint32]func(*AppResponse))
	c.deviceCallbacks = make(map[uint32]*Device)
}

func (c *Client) Address() string {
	if c.useProxy {
		return fmt.Sprintf("wss://companion-rust.facepunch.com/game/%s/%d", c.ip, c.port)
	} else {
		return fmt.Sprintf("ws://%s:%d", c.ip, c.port)
	}
}

func (c *Client) AddToken(token PlayerToken) {
	c.tokens = append(c.tokens, token)
}

func (c *Client) AddDevice(device *Device) {
	c.devices[device.Id] = device
	if c.connection != nil {
		c.initDevice(device)
	}
}

func (c *Client) RemoveDevice(device Device) {
	delete(c.devices, device.Id)
}

func (c *Client) initDevice(device *Device) {
	fmt.Println("connecting device", device.Name)
	c.GetDeviceState(device)
}

func (c *Client) newRequest() (*AppRequest, error) {
	if c.connection == nil {
		return nil, errors.New("connection is nil")
	}
	if len(c.tokens) == 0 {
		return nil, errors.New("no tokens")
	}
	token := c.tokens[0]

	seq := c.seq
	c.seq++
	request := AppRequest{
		Seq:         &seq,
		PlayerId:    &token.SteamId,
		PlayerToken: &token.Token,
	}
	return &request, nil
}

func (c *Client) Write(request *AppRequest) error {
	data, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	c.connection.WriteMessage(websocket.BinaryMessage, data)
	return nil
}

func (c *Client) Read() (*AppResponse, error) {
	_, message, err := c.connection.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("connection error: %s", err)
	}

	if message != nil {
		response := AppMessage{}
		err = proto.Unmarshal(message, &response)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling error: %s", err)
		}

		// If response, execute callback
		if response.Response != nil {
			// Parse any useful data
			c.parseMessageData(response.Response)
			c.executeCallback(response.Response)
			return response.Response, nil
		}
		// If broadcast, handle any updates
		if response.Broadcast != nil {
			c.handleBroadcast(response.Broadcast)
		}
	}

	return nil, nil
}

func (c *Client) parseMessageData(response *AppResponse) {
	if response.EntityInfo != nil {
		if device, ok := c.deviceCallbacks[*response.Seq]; ok {
			device.entType = *response.EntityInfo.Type
		}
	}
}

func (c *Client) handleBroadcast(broadcast *AppBroadcast) {
	if broadcast.EntityChanged != nil {
		if device, ok := c.devices[*broadcast.EntityChanged.EntityId]; ok {
			if device.entType != AppEntityType_StorageMonitor {
				device.value = *broadcast.EntityChanged.Payload.Value
				device.OnUpdate(device)
			} else {
				if *broadcast.EntityChanged.Payload.Value {
					device.items = broadcast.EntityChanged.Payload.Items
					device.capacity = broadcast.EntityChanged.Payload.Capacity
					device.OnUpdate(device)
				}
			}
		}
	}
}

func (c *Client) executeCallback(response *AppResponse) {
	if response.Seq == nil {
		return
	}
	seq := *response.Seq

	if response.EntityInfo == nil {
		// Standard callback pattern
		if callback, ok := c.callbacks[seq]; ok {
			callback(response)
			delete(c.callbacks, seq)
		}
	} else {
		// Handle device callback
		if device, ok := c.deviceCallbacks[seq]; ok {
			if device.OnUpdate != nil {
				device.OnUpdate(device)
			}
			delete(c.deviceCallbacks, seq)
		}
	}
}

func (c *Client) sendRequest(request *AppRequest, callback func(*AppResponse)) error {
	err := c.Write(request)
	if err != nil {
		return err
	}
	if callback != nil {
		c.callbacks[*request.Seq] = callback
	}
	return nil
}
