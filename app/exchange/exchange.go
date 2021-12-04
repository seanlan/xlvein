package exchange

import "context"

type MessageConsume func(message Message)

// Exchange 消息交换器
type Exchange interface {
	Push(message Message)           // 发送消息
	Receive(consume MessageConsume) //消费消息
	Start(ctx context.Context)      // 启动
	Stop()                          // 停止消息接收
}
