package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	WsExchangeName = "xlvein.im"
	ExchangeType   = "fanout"
)

type Session struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	queue amqp.Queue
	Uri   string
}

func NewSession(uri string) *Session {
	return &Session{Uri: uri}
}

func (s *Session) Dial() (err error) {
	var (
		conn  *amqp.Connection
		ch    *amqp.Channel
		queue amqp.Queue
	)
	conn, err = amqp.Dial(s.Uri)
	if err != nil {
		return
	}
	// 建立通道
	ch, err = conn.Channel()
	if err != nil {
		return
	}
	//声明Exchange
	err = ch.ExchangeDeclare(
		WsExchangeName, // name
		ExchangeType,   // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return
	}
	// 声明队列
	u4 := uuid.NewV4()
	queueName := fmt.Sprintf("wool.queue.%s", u4.String())
	queue, err = ch.QueueDeclare(queueName, false, true, true, false, nil)
	if err != nil {
		return
	}
	//绑定队列到Exchange
	err = ch.QueueBind(queueName, "", WsExchangeName, false, nil)
	if err != nil {
		return
	}
	s.Conn = conn
	s.Ch = ch
	s.queue = queue
	return
}

// Close 关闭连接
func (s *Session) Close()  {
	_ = s.Ch.Close()
	_ = s.Conn.Close()
}

func (s *Session) Receive(mq chan Message) (err error) {
	var (
		deliveries <-chan amqp.Delivery
		delivery   amqp.Delivery
	)
	defer func() {
		s.Close()
		zap.S().Info("receive message done")
	}()
	deliveries, err = s.Ch.Consume(
		s.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}
	for {
		select {
		case delivery = <-deliveries:
			if delivery.Body == nil {
				zap.S().Info("receive message", zap.String("body", string(delivery.Body)))
				break
			}
			var msg Message
			err = json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				zap.S().Info("unmarshal message error", zap.Error(err))
				continue
			}
			mq <- msg
		}
	}
}

type RabbitMQExchange struct {
	logger    Logger
	consume   MessageConsume
	messageCh chan Message
	session   *Session
}

func NewRabbitMQExchange(uri string, log Logger) (ex *RabbitMQExchange, err error) {
	ex = &RabbitMQExchange{
		logger: log,
		messageCh: make(chan Message),
	}
	session := NewSession(uri)
	ex.session = session
	err = ex.Dial(uri)
	if err != nil {
		return nil, err
	}
	return
}

func (l *RabbitMQExchange) Dial(uri string) (err error) {
	err = l.session.Dial()
	if err != nil {
		return
	}
	go func() {
		_err := l.session.Receive(l.messageCh)
		if _err != nil {
			return
		}
	}()
	return
}

func (l *RabbitMQExchange) Push(message Message) {
	defer func() {
		if err := recover(); err != nil {
			l.logger.Error("LocalExchange.Push panic: %v", err)
		}
	}()
	var (
		err  error
		body []byte
		msg  amqp.Publishing
	)
	body, err = json.Marshal(message)
	if err != nil {
		return
	}
	msg.Body = body
	msg.DeliveryMode = 2
	err = l.session.Ch.Publish(WsExchangeName, "", false, false, msg)
	if err != nil {
		if e, ok := err.(*amqp.Error); ok {
			if e.Code == amqp.ErrClosed.Code {
				l.logger.Error("LocalExchange.Push reconnect")
				_ = l.Dial(l.session.Uri)
			}
		}
		l.logger.Errorf("LocalExchange.Push error: %v", err)
		return
	}
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
		case msg := <-l.messageCh:
			l.consume(msg)
		}
	}
	l.logger.Debug("LocalExchange.loop exit")
}

func (l *RabbitMQExchange) Start(ctx context.Context) {
	go l.loop(ctx)
}

func (l *RabbitMQExchange) Stop() {
	l.session.Close()
}
