package pkg

type Callback interface {
	Run(c *Client, m *AppResponse)
}
