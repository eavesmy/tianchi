package tianchi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id               string
	conn             *websocket.Conn
	remoteAddr       string
	reConnect        int
	bufferSize       int // unuse
	onMessage        chan []byte
	lastMsgDate      time.Time
	header           http.Header
	sendMessage      chan interface{}
	sendMessageBytes chan []byte
}

// 只有客户端主动调用一次
func (c *Client) connect() (err error) {

	c.conn, _, err = websocket.DefaultDialer.Dial(c.remoteAddr, c.header)

	go c.listen()
	go c.write()

	return
}

func (c *Client) write() {

	var err error

	for {
		select {
		case msg := <-c.sendMessage:
			err = c.conn.WriteJSON(msg)
			break
		case bs := <-c.sendMessageBytes:
			err = c.conn.WriteMessage(websocket.TextMessage, bs)
			break
		}

		if err != nil && pool.hasClient(c.id) {
			pool.errCatch <- err
		}

	}
}

// listen 是一个阻塞操作，持续读取 io.Reader
func (c *Client) listen() {

	for {
		_, bs, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		c.onMessage <- bs
	}
}

func (c *Client) OnMessage() chan []byte {
	return c.onMessage
}

func (c *Client) Send(data interface{}) {
	c.sendMessage <- data
}

func (c *Client) SendBytes(data []byte) {
	c.sendMessageBytes <- data
}

func (c *Client) Handle(handler func(data []byte)) {

	for {
		handler(<-c.onMessage)
	}
}
