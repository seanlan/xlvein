package exchange

type Fields map[string]interface{} // 维护动态字段

type Event int

const (
	EventChatSingle Event = iota + 100 // 单聊消息
	EventChatRoom                      // 房间消息

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
