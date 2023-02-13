package exchange

import "context"

type Fields map[string]interface{} // 维护动态字段

type Event int

const (
	EventSystem   Event = 0 // 系统消息
	EventSingle   Event = 1 // 单聊消息
	EventChatRoom Event = 2 // 群聊消息
)

type Message struct {
	AppID          string `json:"app_id,required"`          // 应用ID
	From           string `json:"from,required"`            //发送者ID
	To             string `json:"to,required"`              //接收者ID
	Event          Event  `json:"event,required"`           //事件类型
	Data           Fields `json:"data,required"`            //消息内容
	MsgID          string `json:"msg_id,required"`          //消息ID
	SendAt         int64  `json:"send_at,required"`         //发送时间
	ConversationID string `json:"conversation_id,required"` //会话ID
}

type MessageConsume func(message Message)

// Exchange 消息交换器
type Exchange interface {
	Push(message Message)           // 将消息推送到交换器
	Receive(consume MessageConsume) // 从交换器接收消息
	Start(ctx context.Context)      // 启动
	Stop()                          // 停止消息接收
}
