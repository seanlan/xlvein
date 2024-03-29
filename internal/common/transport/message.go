package transport

import (
	"github.com/seanlan/xlvein/internal/common/exchange"
)

type Message struct {
	From           string          `json:"from,omitempty"`            //发送者ID
	To             string          `json:"to,required"`               //接收者ID
	Event          exchange.Event  `json:"event,omitempty"`           //事件类型
	Data           exchange.Fields `json:"data,required"`             //消息内容
	MsgID          string          `json:"msg_id,omitempty"`          //消息ID
	SendAt         int64           `json:"send_at,omitempty"`         //发送时间
	ConversationID string          `json:"conversation_id,omitempty"` //会话ID
}
