package tianchi

import "net/http"

type Config struct {
	RemoteAddr     string
	ReConnect      int
	BufferSize     int
	SendBufferSize int
	HttpHeader     http.Header
}

// 初始化
func (c *Config) Init() {
	if c.ReConnect == 0 {
		c.ReConnect = 5
	}

	if c.BufferSize == 0 {
		c.BufferSize = 2048
	}

	if c.SendBufferSize == 0 {
		c.SendBufferSize = 2048
	}
}
