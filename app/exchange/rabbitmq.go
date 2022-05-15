package exchange

import (
	"context"
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

const (
	WsExchangeName = "xlvein.im"
	WxQueueName    = "xlvein.im.queue"
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
	queueName := WxQueueName + "." + u4.String()
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
func (s *Session) Close() {
	_ = s.Ch.Close()
	_ = s.Conn.Close()
}

func (s *Session) Receive(mq chan ExchangeMessage) (err error) {
	var (
		deliveries <-chan amqp.Delivery
		delivery   amqp.Delivery
	)
	defer func() {
		s.Close()
	}()
	deliveries, err = s.Ch.Consume(
		s.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}
	for {
		var loseConnection bool
		select {
		case delivery = <-deliveries:
			if delivery.Body == nil {
				loseConnection = true
				break
			}
			var msg ExchangeMessage
			err = json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				continue
			}
			mq <- msg
		}
		if loseConnection {
			zap.S().Info("rabbitmq lose connection")
			err = errors.New("rabbitmq lose connection")
			break
		}
	}
	return
}

type RabbitMQExchange struct {
	logger    Logger
	consume   MessageConsume
	messageCh chan ExchangeMessage
	session   *Session
}

func NewRabbitMQExchange(uri string, log Logger) (ex *RabbitMQExchange, err error) {
	ex = &RabbitMQExchange{
		logger:    log,
		messageCh: make(chan ExchangeMessage),
	}
	session := NewSession(uri)
	ex.session = session
	err = ex.Dial()
	if err != nil {
		return nil, err
	}
	return
}

func (l *RabbitMQExchange) Dial() (err error) {
	err = l.session.Dial()
	if err != nil {
		return
	}
	go func() {
		for {
			err = l.session.Dial()
			if err != nil {
				l.logger.Error(err)
			}
			err = l.session.Receive(l.messageCh)
			if err != nil {
				l.logger.Error(err)
			}
			time.Sleep(time.Second * 10) // 每5秒重连
		}
	}()
	return
}

func (l *RabbitMQExchange) Push(message ExchangeMessage) {
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
		l.logger.Error("publish message error", zap.Error(err))
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
