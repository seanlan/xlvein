package exchange

import (
	"context"
	"github.com/seanlan/xlvein/internal/common"
)

type LocalExchange struct {
	logger    common.Logger
	consume   MessageConsume
	messageCh chan Message
}

func NewLocalExchange(log common.Logger) (*LocalExchange, error) {
	return &LocalExchange{
		logger:    log,
		messageCh: make(chan Message),
	}, nil
}

// Push 将消息推送到交换器
func (l *LocalExchange) Push(message Message) {
	defer func() {
		if err := recover(); err != nil {
			l.logger.Error("LocalExchange.Push panic: %v", err)
		}
	}()
	l.messageCh <- message
	return
}

// Receive 消费消息（从交换器接收消息）
func (l *LocalExchange) Receive(consume MessageConsume) {
	l.consume = consume
	return
}

func (l *LocalExchange) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-l.messageCh:
			l.consume(message)
		}
	}
}

func (l *LocalExchange) Start(ctx context.Context) {
	go l.loop(ctx)
}

func (l *LocalExchange) Stop() {
	close(l.messageCh)
}
