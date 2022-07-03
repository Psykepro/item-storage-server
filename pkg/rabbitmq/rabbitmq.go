package rabbitmq

import (
	"errors"
	"fmt"

	domain "github.com/Psykepro/item-storage-server/_domain"

	"github.com/Psykepro/item-storage-server/config"

	"github.com/streadway/amqp"
)

// NewRabbitMQConn Initializing new RabbitMQ connection
func NewRabbitMQConn(cfg *config.RabbitMQ) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	return amqp.Dial(connAddr)
}

func InitRabbitMQConn(
	amqpConn *amqp.Connection,
	rabbitMqCfg *config.RabbitMQ,
	logger domain.Logger,
) (*amqp.Channel, <-chan amqp.Delivery, error) {

	ch, err := openRabbitMqChannel(amqpConn, logger)
	if err != nil {
		return nil, nil, err
	}

	q, err := declareRabbitMqQueue(ch, rabbitMqCfg, logger)
	if err != nil {
		return ch, nil, err
	}

	msgs, err := registerRabbitMqConsumer(ch, q, rabbitMqCfg, logger)
	if err != nil {
		return ch, nil, err
	}

	return ch, msgs, nil
}

func openRabbitMqChannel(amqpConn *amqp.Connection, logger domain.Logger) (*amqp.Channel, error) {
	logger.Debugf("Opening RabbitMQ channel ...")
	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open a channel. Err: [%s]", err))
	}
	logger.Debugf("Successfully opened RabbitMQ channel!")

	return ch, nil
}

func registerRabbitMqConsumer(
	ch *amqp.Channel,
	q *amqp.Queue,
	rabbitMqCfg *config.RabbitMQ,
	logger domain.Logger,
) (
	<-chan amqp.Delivery,
	error,
) {
	logger.Debugf("Registering RabbitMQ Consumer ...")
	msgsChannel, err := ch.Consume(
		q.Name,
		rabbitMqCfg.Consumer.Tag,
		rabbitMqCfg.Consumer.AutoAck,
		rabbitMqCfg.Queue.Exclusive,
		rabbitMqCfg.Consumer.NoLocal,
		rabbitMqCfg.Queue.NoWait,
		rabbitMqCfg.Consumer.Args,
	)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to register RabbitMQ Consumer. Err: [%s]", err))
	}
	logger.Debugf("Successfully registered RabbitMQ Consumer!")

	return msgsChannel, nil
}

func declareRabbitMqQueue(ch *amqp.Channel, rabbitMqCfg *config.RabbitMQ, logger domain.Logger) (*amqp.Queue, error) {
	logger.Debugf("Declaring RabbitMQ Queue ...")
	q, err := ch.QueueDeclare(
		rabbitMqCfg.Queue.Name,
		rabbitMqCfg.Queue.Durable,
		rabbitMqCfg.Queue.AutoDelete,
		rabbitMqCfg.Queue.Exclusive,
		rabbitMqCfg.Queue.NoWait,
		rabbitMqCfg.Queue.Args,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to declare RabbitMQ Queue. Err: [%s]", err))
	}
	logger.Debugf("Successfully declaring RabbitMQ Queue!")

	return &q, nil
}
