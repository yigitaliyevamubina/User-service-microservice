package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"template-service3/service"

	pb "template-service3/genproto/user_service"

	"github.com/k0kubun/pp"
	"github.com/streadway/amqp"
)

type consumer interface {
	ConsumeMessages(queue string, handler func([]byte, *service.UserService)) error
	Close() error
}

type RabbitMQConsumer struct {
	channel *amqp.Channel
	service *service.UserService
}

func NewRabbitMQConsumer(channel *amqp.Channel, service *service.UserService) consumer {
	return &RabbitMQConsumer{channel: channel, service: service}
}

func (c *RabbitMQConsumer) ConsumeMessages(queue string, handler func([]byte, *service.UserService)) error {
	pp.Println("consuming from queue...")
	q, err := c.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = c.channel.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	if err != nil {
		return err
	}
	messages, err := c.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			handler(msg.Body, c.service)
		}
	}()

	return nil
}

func (c *RabbitMQConsumer) Close() error {
	return c.channel.Close()
}

func ConsumeHandler(message []byte, service *service.UserService) {
	var user pb.User
	if err := json.Unmarshal(message, &user); err != nil {
		log.Println(err)
		return
	}

	respUser, err := service.Create(context.Background(), &user)
	if err != nil {
		log.Fatal("cannot create user via rabbitmq")
	}

	pp.Println(respUser)
}
