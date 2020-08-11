/*
	不适用大型项目。
*/
package tianchi

import (
	"github.com/satori/go.uuid"
	"net/http"
)

var pool *Pool

// 暴露获取 client 接口

// 建立链接池

// 由协议升级
func NewClient(remoteAddr string, confs ...Config) (c *Client, err error) {
	conf := &Config{
		RemoteAddr:     remoteAddr,
		BufferSize:     2048,
		ReConnect:      5,
		HttpHeader:     http.Header{},
		SendBufferSize: 2048,
	}
	uid, _ := uuid.NewV4()
	c = &Client{
		id:               uid.String(),
		remoteAddr:       conf.RemoteAddr,
		onMessage:        make(chan []byte, conf.BufferSize),
		reConnect:        conf.ReConnect,
		header:           conf.HttpHeader,
		sendMessage:      make(chan interface{}, conf.SendBufferSize),
		sendMessageBytes: make(chan []byte, conf.SendBufferSize),
	}
	err = c.connect()
	return
}

func Register(c *Client) {
	pool.Register(c)
}

type Pool struct {
	cs       []*Client
	errCatch chan error
}

func (p *Pool) Register(c *Client) {
	if len(p.cs) == 0 {
		p.cs = append(p.cs, c)
		return
	}

	index := -1
	for i, c := range p.cs {
		if c == nil {
			index = i
		}
	}

	if index == -1 {
		p.cs = append(p.cs, c)
	}
}

func (p *Pool) hasClient(id string) (has bool) {
	for _, c := range p.cs {
		if c.id == id {
			return true
		}
	}
	return
}

func (p *Pool) CatchErr(handle func(error)) {
	for {
		handle(<-p.errCatch)
	}
}

func New() (pool *Pool) {
	pool = &Pool{
		cs:       []*Client{},
		errCatch: make(chan error, 10),
	}
	return
}
