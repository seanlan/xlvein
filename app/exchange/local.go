package exchange

import (
	"context"
)

type LocalExchange struct {
	logger    Logger
	consume   MessageConsume
	messageCh chan Message
}

func NewLocalExchange(l Logger) *LocalExchange {
	return &LocalExchange{
		messageCh: make(chan Message, 100),
	}
}

func (l *LocalExchange) Push(message Message) {
	defer func() {
		if err := recover(); err != nil {
			l.logger.Error("LocalExchange.Push panic: %v", err)
		}
	}()
	l.messageCh <- message
	return
}

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
