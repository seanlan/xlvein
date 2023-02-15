package veinsdk

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/seanlan/xlvein/internal/common/transport"
)

type ReceiveHandler func(msg transport.Message)

type Client struct {
	Token   string
	Uri     string
	Cli     *websocket.Conn
	Handler ReceiveHandler // 接收消息处理函数
}

//
// NewClient 创建一个新的客户端
//  @Description: NewClient 创建一个新的客户端
//  @param token
//  @param uri
//  @return *Client
//
func NewClient(uri string, handler ReceiveHandler) *Client {
	return &Client{
		Uri:     uri,
		Handler: handler,
	}
}

// Connect 建立连接
func (c *Client) connect() (err error) {
	c.Cli, _, err = websocket.DefaultDialer.Dial(c.Uri, nil)
	return
}

// Close 关闭连接
func (c *Client) close() {
	c.Cli.Close()
}

// Send 发送消息
func (c *Client) Send(msg transport.Message) (err error) {
	var buff []byte
	buff, err = json.Marshal(msg)
	if err != nil {
		return
	}
	err = c.Cli.WriteMessage(websocket.TextMessage, buff)
	return
}

// Receive 接收消息
func (c *Client) Receive() (err error) {
	var (
		msgType int
		buff    []byte
		msg     transport.Message
	)
	msgType, buff, err = c.Cli.ReadMessage()
	if err != nil {
		return
	}
	if msgType == websocket.TextMessage {
		err = json.Unmarshal(buff, &msg)
		if err != nil {
			return
		}
		c.Handler(msg)
	}
	return
}

// Run 运行客户端
func (c *Client) Run() (err error) {
	err = c.connect()
	if err != nil {
		return
	}
	defer c.close()
	for {
		err = c.Receive()
		if err != nil {
			return
		}
	}
}
