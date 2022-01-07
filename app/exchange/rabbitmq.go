package exchange

import "context"

type RabbitMQExchange struct {
	logger    Logger
	consume   MessageConsume
	messageCh chan Message
}

func NewRabbitMQExchange(log Logger) *RabbitMQExchange {
	return &RabbitMQExchange{
		logger:    log,
		messageCh: make(chan Message),
	}
}

func (l *RabbitMQExchange) Push(message Message) {
	defer func() {
		if err := recover(); err != nil {
			l.logger.Error("LocalExchange.Push panic: %v", err)
		}
	}()
	l.messageCh <- message
	return
}

func (l *RabbitMQExchange) Receive(consume MessageConsume) {
	l.consume = consume
	return
}

func (l *RabbitMQExchange) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-l.messageCh:
			l.consume(message)
		}
	}
}

func (l *RabbitMQExchange) Start(ctx context.Context) {
	go l.loop(ctx)
}

func (l *RabbitMQExchange) Stop() {
	close(l.messageCh)
}

