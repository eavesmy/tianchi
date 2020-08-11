package tianchi

import "net/http"

type Config struct {
	RemoteAddr     string
	ReConnect      int
	BufferSize     int
	SendBufferSize int
	HttpHeader     http.Header
}
