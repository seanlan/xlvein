package exchange

import (
	"context"
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/seanlan/xlvein/app/common"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

type Session struct {
	ExchangeName string
	QueueName    string
	ExchangeType string
	Conn         *amqp.Connection
	Ch           *amqp.Channel
	queue        amqp.Queue
	Uri          string
}

func NewSession(uri, exchangeName, queueName string) *Session {
	return &Session{
		Uri:          uri,
		ExchangeName: exchangeName,
		QueueName:    queueName,
		ExchangeType: "fanout",
	}
}

func (s *Session) Dial() (err error) { // 连接rabbitmq
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
		s.ExchangeName, // name
		s.ExchangeType, // type
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
	queueName := s.QueueName + "." + u4.String()
	queue, err = ch.QueueDeclare(queueName, false, true, true, false, nil)
	if err != nil {
		return
	}
	//绑定队列到Exchange
	err = ch.QueueBind(queueName, "", s.ExchangeName, false, nil)
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

func (s *Session) Receive(mq chan Message) (err error) {
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
			var msg Message
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
	logger    common.Logger
	consume   MessageConsume
	messageCh chan Message
	session   *Session
}

func NewRabbitMQExchange(uri, exchangeName, queueName string, log common.Logger) (ex *RabbitMQExchange, err error) {
	ex = &RabbitMQExchange{
		logger:    log,
		messageCh: make(chan Message),
	}
	session := NewSession(uri, exchangeName, queueName)
	ex.session = session
	err = ex.Dial()
	if err != nil {
		return nil, err
	}
	return
}

func (ex *RabbitMQExchange) Dial() (err error) {
	err = ex.session.Dial()
	if err != nil {
		ex.logger.Error(err)
		return
	}
	go func() {
		for {
			err = ex.session.Dial()
			if err != nil {
				ex.logger.Error(err)
			}
			err = ex.session.Receive(ex.messageCh)
			if err != nil {
				ex.logger.Error(err)
			}
			time.Sleep(time.Second * 10) // 每10秒重连
		}
	}()
	return
}

func (ex *RabbitMQExchange) Push(message Message) {
	defer func() {
		if err := recover(); err != nil {
			ex.logger.Error("LocalExchange.Push panic: %v", err)
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
	err = ex.session.Ch.Publish(ex.session.ExchangeName, "", false, false, msg)
	if err != nil {
		ex.logger.Error("publish message error", zap.Error(err))
		return
	}
	return
}

func (ex *RabbitMQExchange) Receive(consume MessageConsume) {
	ex.consume = consume
	return
}

func (ex *RabbitMQExchange) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ex.messageCh:
			ex.consume(msg)
		}
	}
	ex.logger.Debug("LocalExchange.loop exit")
}

func (ex *RabbitMQExchange) Start(ctx context.Context) {
	go ex.loop(ctx)
}

func (ex *RabbitMQExchange) Stop() {
	ex.session.Close()
}
