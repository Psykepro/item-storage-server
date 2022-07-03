package server

import (
	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
	domain "github.com/Psykepro/item-storage-server/_domain"
	"github.com/Psykepro/item-storage-server/config"
	"github.com/Psykepro/item-storage-server/pkg/rabbitmq"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	rabbitMqCfg        *config.RabbitMQ
	itemRequestHandler domain.ItemRequestHandler
	logger             domain.Logger
}

func NewServer(rabbitMqCfg *config.RabbitMQ, itemRequestHandler domain.ItemRequestHandler, logger domain.Logger) *Server {
	return &Server{
		rabbitMqCfg:        rabbitMqCfg,
		itemRequestHandler: itemRequestHandler,
		logger:             logger,
	}
}

func (s *Server) Run() {

	amqpConn, rabbitMqCh, msgChannel, err := s.initToRabbitMq()
	if err != nil {
		s.logger.Fatalf(err.Error())
	}
	defer amqpConn.Close()
	defer rabbitMqCh.Close()

	// Starting handler
	itemRequestChannel := make(chan *pb.ItemRequest)
	go s.itemRequestHandler.Handle(itemRequestChannel)

	// Endless Consuming
	for {
		for d := range msgChannel {
			request := &pb.ItemRequest{}
			if err = proto.Unmarshal(d.Body, request); err != nil {
				s.logger.Errorf("Failed to parse Item Request. Err: [%s]", err)
			}
			s.logger.Debugf("Received a request message for command [%s]", request.Command)
			itemRequestChannel <- request
		}
		s.logger.Debugf(" [*] Waiting for messages.")
	}
}

func (s *Server) initToRabbitMq() (*amqp.Connection, *amqp.Channel, <-chan amqp.Delivery, error) {
	s.logger.Debugf("Connecting to RabbitMQ ...")
	amqpConn, err := rabbitmq.NewRabbitMQConn(s.rabbitMqCfg)
	if err != nil {
		s.fatalOnError("Failed to connect to RabbitMQ", err)
	}
	s.logger.Debugf("Successfully connected to RabbitMQ!")
	rabbitMqChannel, msgChannel, err := rabbitmq.InitRabbitMQConn(amqpConn, s.rabbitMqCfg, s.logger)

	return amqpConn, rabbitMqChannel, msgChannel, err
}

func (s *Server) fatalOnError(msg string, err error) {
	if err != nil {
		s.logger.Fatalf("%s. Err: [%s]", msg, err)
	}
}
