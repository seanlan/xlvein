package transport

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/seanlan/xlvein/app/common"
	"github.com/seanlan/xlvein/app/common/exchange"
)

type GetConversationMembers func(appID, conversationId string) ([]string, error) // 获取会话成员

type Hub struct { // transport 管理器
	exchange      exchange.Exchange
	clients       map[string]*hashset.Set
	logger        common.Logger
	memberHandler GetConversationMembers
}

var ClientHub *Hub

func InitHub(ctx context.Context, ex exchange.Exchange, logger common.Logger) {
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
func (h *Hub) PushToExchange(appID string, msg Message) {
	uu, _ := uuid.NewUUID()
	var exchangeMsg = exchange.Message{
		AppID:          appID,
		From:           msg.From,
		To:             msg.To,
		Event:          msg.Event,
		Data:           msg.Data,
		MsgID:          uu.String(),
		ConversationID: msg.ConversationID,
	}
	h.exchange.Push(exchangeMsg)
	//TODO 可以在这里记录消息历史记录
}

// 发送到指定的客户端
func (h *Hub) sendToTransport(appID, to string, msg exchange.Message) {
	key := makeTransportKey(appID, to)
	m := Message{
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
func (h *Hub) distribute(message exchange.Message) {
	switch message.Event {
	case exchange.EventSystem, exchange.EventSingle: // 系统消息||单聊消息，直接发送
		h.sendToTransport(message.AppID, message.To, message)
	case exchange.EventChatRoom: // 群聊
		// 需要找到会话中的所有成员，发送消息
		members, err := h.memberHandler(message.AppID, message.ConversationID)
		if err != nil {
			return
		}
		for _, member := range members {
			h.sendToTransport(message.AppID, member, message)
		}
	}
}

// Run 启动
func (h *Hub) Run(ctx context.Context) {
	h.exchange.Receive(func(message exchange.Message) {
		h.distribute(message)
	})
	h.exchange.Start(ctx)
}
