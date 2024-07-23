package exchange

import "context"

const (
	TypeLocal    = "local"
	TypeRabbitMQ = "rabbitmq"
	TypeRedis    = "redis"
)

// Exchange 消息交换器
type Exchange interface {
	Push(message Message)           // 将消息推送到交换器
	Receive(consume MessageConsume) // 从交换器接收消息
	Run(ctx context.Context)        // 启动
	Stop()                          // 停止消息接收
}
