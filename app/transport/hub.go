package transport

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/seanlan/xlvein/app/exchange"
	"sort"
	"strings"
	"time"
)

// MakeConversationID 构建会话ID
func MakeConversationID(appID, from, to string, event exchange.Event) string {
	if event == exchange.EventChatSingle {
		keys := []string{from, to}
		sort.Strings(keys)
		source := strings.Join(keys, ":")
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%s:%s", appID, source)))
		return hex.EncodeToString(h.Sum(nil))
	} else {
		return to
	}
}

type Hub struct {
	Exchange exchange.Exchange
	// clients Map
	clients map[string]*hashset.Set
	logger  exchange.Logger
}

var ClientHub *Hub

func InitHub(ctx context.Context, ex exchange.Exchange, logger exchange.Logger) {
	ClientHub = &Hub{
		Exchange: ex,
		clients:  make(map[string]*hashset.Set),
		logger:   logger,
	}
	ClientHub.Run(ctx)
}

// Join 加入一条链接
func (h *Hub) Join(appID, Tag string, conn *websocket.Conn) {
	transport := NewTransport(appID, Tag, conn, h)
	transport.Start()
	if h.clients[transport.Key] == nil {
		h.clients[transport.Key] = hashset.New()
	}
	h.clients[transport.Key].Add(transport)
}

// Drop 剔除一条链接
func (h *Hub) Drop(transport *Transport) {
	key := transport.Key
	if _, ok := h.clients[key]; ok {
		h.clients[key].Remove(transport)
	}
}

// PushToExchange 将消费推送到消息交换器
func (h *Hub) PushToExchange(appID string, msg Message) {
	uu, _ := uuid.NewUUID()
	var exchangeMsg = exchange.Message{
		AppID:          appID,
		From:           msg.From,
		To:             msg.To,
		Event:          msg.Event,
		Data:           msg.Data,
		MsgID:          uu.String(),
		SendAt:         time.Now().Unix(),
		ConversationID: MakeConversationID(appID, msg.From, msg.To, msg.Event),
	}
	h.Exchange.Push(exchangeMsg)
}

// 发送到指定的客户端
func (h *Hub) sendToTransport(msg exchange.Message) {
	key := makeTransportKey(msg.AppID, msg.To)
	m := Message{
		From:           msg.From,
		To:             msg.To,
		Event:          msg.Event,
		Data:           msg.Data,
		MsgID:          msg.MsgID,
		SendAt:         msg.SendAt,
		ConversationID: msg.ConversationID,
	}
	h.logger.Debugf("sendToTransport: %s", key)
	h.logger.Debugf("sendToTransport: %+v", m)
	if sets, ok := h.clients[key]; ok {
		for _, t := range sets.Values() {
			t.(*Transport).Send(m)
		}
	}
}

// 消息分发
func (h *Hub) distribute(message exchange.Message) {
	switch message.Event {
	case exchange.EventChatSingle:
		h.sendToTransport(message)
	case exchange.EventChatRoom:
		h.sendToTransport(message)
	}
}

// Run 启动
func (h *Hub) Run(ctx context.Context) {
	h.Exchange.Receive(func(message exchange.Message) {
		h.distribute(message)
	})
	h.Exchange.Start(ctx)
}
