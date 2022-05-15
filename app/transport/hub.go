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
)

// MakeConversationID 构建会话ID
func MakeConversationID(appID, from, to string, event exchange.ExchangeEvent) string {
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
	exchange exchange.Exchange
	// clients Map
	clients map[string]*hashset.Set
	logger  exchange.Logger
}

var ClientHub *Hub

func InitHub(ctx context.Context, ex exchange.Exchange, logger exchange.Logger) {
	ClientHub = &Hub{
		exchange: ex,
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
func (h *Hub) PushToExchange(appID string, msg TransportMessage) {
	uu, _ := uuid.NewUUID()
	var exchangeMsg = exchange.ExchangeMessage{
		AppID:          appID,
		From:           msg.From,
		To:             msg.To,
		Event:          msg.Event,
		Data:           msg.Data,
		MsgID:          uu.String(),
		ConversationID: MakeConversationID(appID, msg.From, msg.To, msg.Event),
	}
	h.exchange.Push(exchangeMsg)
}

// 发送到指定的客户端
func (h *Hub) sendToTransport(msg exchange.ExchangeMessage) {
	key := makeTransportKey(msg.AppID, msg.To)
	m := TransportMessage{
		From:           msg.From,
		To:             msg.To,
		Event:          msg.Event,
		Data:           msg.Data,
		MsgID:          msg.MsgID,
		SendAt:         msg.SendAt,
		ConversationID: msg.ConversationID,
	}
	if sets, ok := h.clients[key]; ok {
		for _, t := range sets.Values() {
			t.(*Transport).Send(m)
		}
	}
}

// 消息分发
func (h *Hub) distribute(message exchange.ExchangeMessage) {
	switch message.Event {
	case exchange.EventChatSingle: // 单聊
		h.sendToTransport(message)
	case exchange.EventChatRoom:  // 群聊
		h.sendToTransport(message)
	}
}

// Run 启动
func (h *Hub) Run(ctx context.Context) {
	h.exchange.Receive(func(message exchange.ExchangeMessage) {
		h.distribute(message)
	})
	h.exchange.Start(ctx)
}
