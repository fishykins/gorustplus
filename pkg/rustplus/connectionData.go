package rustplus

import "fmt"

type ConnectionData struct {
	useProxy bool
	ip       string
	port     uint64
	Tokens   []PlayerToken
}

func NewConnectionData(ip string, port uint64, useProxy bool) ConnectionData {
	return ConnectionData{
		ip:       ip,
		port:     port,
		useProxy: useProxy,
		Tokens:   make([]PlayerToken, 0),
	}
}

func (c *ConnectionData) AddToken(token PlayerToken) {
	c.Tokens = append(c.Tokens, token)
}

func (c *ConnectionData) URL() string {
	if c.useProxy {
		return fmt.Sprintf("wss://companion-rust.facepunch.com/game/%s/%d", c.ip, c.port)
	} else {
		return fmt.Sprintf("ws://%s:%d", c.ip, c.port)
	}
}
