package communication

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMQCommunicator struct {
	ctx              context.Context
	consumerConn     *amqp.Connection
	publisherConn    *amqp.Connection
	reconnectTimeout time.Duration
	toSendBuf        [][]byte
	connStr          string
}

func InitRabbitMQCommunicator(
	ctx context.Context,
	connStr string,
	reconnectTimeout time.Duration,
) (*RabbitMQCommunicator, error) {
	consumerConn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	publisherConn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	return &RabbitMQCommunicator{
		ctx:              ctx,
		consumerConn:     consumerConn,
		publisherConn:    publisherConn,
		reconnectTimeout: reconnectTimeout,
		toSendBuf:        make([][]byte, 0),
		connStr:          connStr,
	}, nil
}

func (comm *RabbitMQCommunicator) Close() {
	_ = comm.consumerConn.Close()
	_ = comm.publisherConn.Close()
}

func (comm *RabbitMQCommunicator) DeclareExchange(exchangeName string) error {
	ch, err := comm.consumerConn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	err = ch.ExchangeDeclare(
		exchangeName,        // exchange name
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // autoDelete
		false,               // internal
		false,               // noWait
		nil,                 // args
	)
	return err
}

func (comm *RabbitMQCommunicator) DeclareQueueAndBind(queueName string, exchangeName string) error {
	ch, err := comm.consumerConn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // nowait
		nil,          // args
	)
	return err
}

func (comm *RabbitMQCommunicator) publisher(
	ctx context.Context,
	exchangeName string,
	dataSource <-chan []byte,
	ch *amqp.Channel,
	logger *log.Logger,
) {
	defer func() {
		_ = ch.Close()
	}()
	counter := 0
	logger.Printf("start")

	if len(comm.toSendBuf) != 0 {
		logger.Printf("there are %v messages pending. Sending...", len(comm.toSendBuf))
		for len(comm.toSendBuf) != 0 {
			err := ch.PublishWithContext(
				ctx,
				exchangeName, // exchange
				"",           // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        comm.toSendBuf[0],
				})
			if err != nil {
				logger.Printf("error while publishing: %s", err)
				var mqErr *amqp.Error
				if errors.As(err, &mqErr) {
					go comm.reconnectPublisher(exchangeName, dataSource)
				}
				return
			}
			comm.toSendBuf = comm.toSendBuf[1:]
			counter += 1
			logger.Printf("published %v messages", counter)
		}
		logger.Printf("all messages sent")
	}

	for {
		select {
		case <-ctx.Done():
			logger.Printf("context.Done")
			return
		case data, ok := <-dataSource:
			if !ok {
				logger.Printf("data source channel closed")
			}
			err := ch.PublishWithContext(
				ctx,
				exchangeName, // exchange
				"",           // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        data,
				})
			if err != nil {
				logger.Printf("error while publishing: %s", err)
				var mqErr *amqp.Error
				if errors.As(err, &mqErr) {
					comm.toSendBuf = append(comm.toSendBuf, data)
					go comm.reconnectPublisher(exchangeName, dataSource)
				}
				return
			}
			counter += 1
			logger.Printf("published %v messages", counter)
		}
	}
}

func (comm *RabbitMQCommunicator) RunPublisher(exchangeName string, dataSource <-chan []byte) error {
	ch, err := comm.publisherConn.Channel()
	if err != nil {
		return err
	}
	defaultLogger := log.Default()
	logger := log.New(
		defaultLogger.Writer(),
		fmt.Sprintf("publisher (%s): ", exchangeName),
		defaultLogger.Flags()|log.Lmsgprefix)
	go comm.publisher(comm.ctx, exchangeName, dataSource, ch, logger)
	return nil
}

func (comm *RabbitMQCommunicator) consumer(
	ctx context.Context,
	queueName string,
	f func(data []byte, logger *log.Logger) error,
	ch *amqp.Channel,
	logger *log.Logger,
) {
	defer func() {
		_ = ch.Close()
	}()
	counter := 0
	logger.Printf("start")
	msgs, err := ch.ConsumeWithContext(
		ctx,
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logger.Printf("error after call to ConsumeWithContext: %s", err)
		var mqErr *amqp.Error
		if errors.As(err, &mqErr) {
			go comm.reconnectConsumer(queueName, f)
		}
		return
	}
	closeInfo := ch.NotifyClose(make(chan *amqp.Error, 1))
	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				logger.Printf("no msgs from channel. exiting...")
				continue
			}
			counter += 1
			logger.Printf("consumed %v messages", counter)
			err := f(d.Body, logger)
			if err != nil {
				logger.Printf("error after executing user function: %s", err)
			}
		case err := <-closeInfo:
			if err != nil {
				go comm.reconnectConsumer(queueName, f)
			}
			return
		}
	}
}

func (comm *RabbitMQCommunicator) RunConsumer(queueName string, f func(data []byte, logger *log.Logger) error) error {
	ch, err := comm.consumerConn.Channel()
	if err != nil {
		return err
	}
	defaultLogger := log.Default()
	logger := log.New(
		defaultLogger.Writer(),
		fmt.Sprintf("consumer (%s): ", queueName),
		defaultLogger.Flags()|log.Lmsgprefix)
	go comm.consumer(comm.ctx, queueName, f, ch, logger)
	return nil
}

func (comm *RabbitMQCommunicator) reconnectPublisher(
	exchangeName string,
	dataSource <-chan []byte,
) {
	_ = comm.publisherConn.Close()
	log.Printf("publisher tries to reconnect ...")
	timer := time.NewTimer(comm.reconnectTimeout)
	for {
		select {
		case <-comm.ctx.Done():
			log.Printf("context.Done")
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			publisherConn, err := amqp.Dial(comm.connStr)
			if err != nil {
				log.Printf("publisher reconnect: failed to establish connection: %s", err)
				timer.Reset(comm.reconnectTimeout)
				continue
			}
			comm.publisherConn = publisherConn
			err = comm.RunPublisher(exchangeName, dataSource)
			if err != nil {
				log.Printf("publisher reconnect: failed to run publisher: %s", err)
				_ = comm.publisherConn.Close()
				timer.Reset(comm.reconnectTimeout)
				continue
			}
			return
		}
	}
}

func (comm *RabbitMQCommunicator) reconnectConsumer(
	queueName string,
	f func(data []byte, logger *log.Logger) error,
) {
	_ = comm.consumerConn.Close()
	log.Printf("consumer tries to reconnect ...")
	timer := time.NewTimer(comm.reconnectTimeout)
	for {
		select {
		case <-comm.ctx.Done():
			log.Printf("context.Done")
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			consumerConn, err := amqp.Dial(comm.connStr)
			if err != nil {
				log.Printf("consumer reconnect: failed to establish connection: %s", err)
				timer.Reset(comm.reconnectTimeout)
				continue
			}
			comm.consumerConn = consumerConn
			err = comm.RunConsumer(queueName, f)
			if err != nil {
				log.Printf("publisher reconnect: failed to run publisher: %s", err)
				_ = comm.consumerConn.Close()
				timer.Reset(comm.reconnectTimeout)
				continue
			}
			return
		}
	}
}
