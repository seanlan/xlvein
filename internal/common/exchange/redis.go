package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/seanlan/xlvein/internal/common"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

type RedisConn struct {
	Uri       string // redis地址
	QueueName string
	Conn      *redis.Client
	Ch        *redis.PubSub
}

func NewRedisConn(uri, queueName string) *RedisConn {
	return &RedisConn{
		Uri:       uri,
		QueueName: queueName,
	}
}

func (c *RedisConn) Dial(ctx context.Context) (err error) {
	c.Conn = redis.NewClient(&redis.Options{
		Addr: c.Uri,
	})
	c.Ch = c.Conn.Subscribe(ctx, c.QueueName)
	return
}

func (c *RedisConn) Close() (err error) {
	return c.Ch.Close()
}

func (c *RedisConn) Receive(ctx context.Context, mq chan Message) (err error) {
	defer func() {
		c.Close()
	}()
	_, err = c.Ch.Receive(ctx)
	if err != nil {
		return
	}
	deliveries := c.Ch.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-deliveries:
			if !ok {
				return errors.New("channel closed")
			}
			var m Message
			err = json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				return
			}
			mq <- m
		}
	}
}

type RedisExchange struct {
	logger    common.Logger
	consume   MessageConsume
	messageCh chan Message
	conn      *RedisConn
}

func NewRedisExchange(ctx context.Context, uri, queueName string, log common.Logger) (ex *RedisExchange, err error) {
	ex = &RedisExchange{
		logger:    log,
		messageCh: make(chan Message),
	}
	conn := NewRedisConn(uri, queueName)
	ex.conn = conn
	err = ex.Dial(ctx)
	return
}

func (ex *RedisExchange) Dial(ctx context.Context) (err error) {
	err = ex.conn.Dial(ctx)
	if err != nil {
		ex.logger.Error(err)
		return
	}
	go func() {
		for {
			err = ex.conn.Dial(ctx)
			if err != nil {
				ex.logger.Error(err)
			}
			err = ex.conn.Receive(ctx, ex.messageCh)
			if err != nil {
				ex.logger.Error(err)
			}
			time.Sleep(time.Second * 10) // 每10秒重连
		}
	}()
	return
}

func (ex *RedisExchange) Push(message Message) {
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
	err = ex.conn.Conn.Publish(context.Background(), ex.conn.QueueName, msg.Body).Err()
	if err != nil {
		ex.logger.Error("publish message error", zap.Error(err))
		return
	}
	return
}

func (ex *RedisExchange) Receive(consume MessageConsume) {
	ex.consume = consume
	return
}

func (ex *RedisExchange) loop(ctx context.Context) {
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

func (ex *RedisExchange) Start(ctx context.Context) {
	go ex.loop(ctx)
}

func (ex *RedisExchange) Stop() {
	ex.conn.Close()
}
