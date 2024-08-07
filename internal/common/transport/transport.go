package transport

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const (
	writeWait  = 10 * time.Second    // 发送消息超时时间 10秒
	pongWait   = 10 * time.Second    // 响应ping消息超时时间 10秒
	pingPeriod = (pongWait * 9) / 10 // 发送ping消息间隔时间 9秒
)

// makeTransportKey 根据appID和tag生成client唯一标示
func makeTransportKey(appID, tag string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", appID, tag)))
	return hex.EncodeToString(h.Sum(nil))
}

type Transport struct {
	AppID    string          // 应用ID
	Tag      string          // 连接标签
	Key      string          // 连接标识  = makeTransportKey(AppID:Tag)
	hub      *Hub            // 会话管理
	conn     *websocket.Conn //websocket连接
	sendPool chan Message    // 发送消息队列
}

// NewTransport 创建一个新的连接
func NewTransport(appID, tag string, conn *websocket.Conn, hub *Hub) *Transport {
	return &Transport{
		AppID:    appID,
		Tag:      tag,
		Key:      makeTransportKey(appID, tag),
		hub:      hub,
		conn:     conn,
		sendPool: make(chan Message, 100),
	}
}

func (trans *Transport) Run() {
	go trans.doRead()
	go trans.doWrite()
}

// 启动消息读取
func (trans *Transport) doRead() {
	defer func() {
		trans.hub.logger.Debugf("trans.hub.clients : %+v", trans.hub.clients)
		trans.hub.Drop(trans)
		_ = trans.conn.Close()
		close(trans.sendPool)
		trans.hub.logger.Debugf("trans.hub.clients : %+v", trans.hub.clients)
		trans.hub.logger.Debugf("doRead defer")
	}()
	_ = trans.conn.SetReadDeadline(time.Now().Add(pongWait))
	// 设置心跳处理
	trans.conn.SetPongHandler(
		func(string) error {
			_ = trans.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	// 开始接收消息
	for {
		messageType, message, err := trans.conn.ReadMessage()
		if err != nil {
			trans.hub.logger.Infof("read message error: %s", err.Error())
			return
		}
		if messageType == websocket.PingMessage {
			trans.hub.logger.Info("this is a  pingMessage")
		}
		if messageType == websocket.CloseMessage {
			trans.hub.logger.Info("this is a  closeMessage")
			return
		}
		trans.hub.logger.Debugf("received messages: %s", string(message))
		//消息重新封装
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			// 消息格式解析失败
			trans.hub.logger.Debugf("Unmarshal failed:%s", message)
			continue
		}
		//msg.Event = msg
		msg.From = trans.Tag
		msg.SendAt = time.Now().Unix()
		msg.ConversationID = msg.To
		// 将消息推送到消息交换器
		trans.hub.PushToExchange(trans.AppID, msg)
	}
	trans.hub.logger.Debug("transport doRead stop")
}

// 启动消息发送
func (trans *Transport) doWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = trans.conn.Close()
		trans.hub.logger.Debugf("doWrite defer trans.conn.Close")
	}()
	for {
		select {
		case message, ok := <-trans.sendPool:
			_ = trans.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = trans.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			var (
				poolLen = len(trans.sendPool)
				pool    = make([]Message, 0, poolLen+1)
			)
			pool = append(pool, message)
			for i := 0; i < poolLen; i++ {
				pool = append(pool, <-trans.sendPool)
			}
			for _, msg := range pool {
				buff, err := json.Marshal(msg)
				if err == nil {
					err = trans.conn.WriteMessage(websocket.TextMessage, buff)
					if err != nil {
						return
					}
				}
			}
		case <-ticker.C:
			_ = trans.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := trans.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send 发送消息
func (trans *Transport) Send(message Message) {
	trans.sendPool <- message
}
