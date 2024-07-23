package exchange

type Fields map[string]interface{} // 维护动态字段

func (f Fields) String(key string) string {
	if v, ok := f[key]; ok {
		return v.(string)
	}
	return ""
}

func (f Fields) Int(key string) int {
	if v, ok := f[key]; ok {
		return v.(int)
	}
	return 0
}

func (f Fields) Int64(key string) int64 {
	if v, ok := f[key]; ok {
		return v.(int64)
	}
	return 0
}

func (f Fields) Float64(key string) float64 {
	if v, ok := f[key]; ok {
		return v.(float64)
	}
	return 0
}

func (f Fields) Bool(key string) bool {
	if v, ok := f[key]; ok {
		return v.(bool)
	}
	return false
}

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
