package messaging

import (
	"github.com/streadway/amqp"
)

type RabbitMQBus struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	url        string
	id         string
}

func NewRabbitMQBus(url string, id string) (*RabbitMQBus, error) {
	var err error
	bus := new(RabbitMQBus)
	bus.url = url
	bus.id = id
	bus.connection, err = amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	bus.channel, err = bus.connection.Channel()
	if err != nil {
		return nil, err
	}
	return bus, nil
}

func (r *RabbitMQBus) CreateRabbitMQWorker() MessageBusConsumer {
	wq := new(RabbitMQWorker)
	wq.channel = r.channel
	wq.id = r.id
	return wq
}

func (r *RabbitMQBus) SendMessageToRabbitMQWorkQueue(msg string) {
	sendMessageToRabbitMQWorkQueue(r.url, r.id, msg)
}

func (r *RabbitMQBus) Close() {
	r.channel.Close()
	r.connection.Close()
}
