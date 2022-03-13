package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQWorker struct {
	channel     *amqp.Channel
	msgReceived func(msg string) error
	id          string
}

func (r *RabbitMQWorker) OnMessageReceived(f func(msg string) error) {
	r.msgReceived = f
}

func (r *RabbitMQWorker) StartListening() error {
	queue, err := r.channel.QueueDeclare(
		fmt.Sprintf("wq-%s", r.id),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msgs, err := r.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go func() {
		for m := range msgs {
			r.msgReceived(string(m.Body)) // check error
			m.Ack(false)
		}
	}()

	return nil
}

func sendMessageToRabbitMQWorkQueue(url string, id string, msg string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer conn.Close()
	queue, err := channel.QueueDeclare(
		fmt.Sprintf("wq-%s", id),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(msg),
		})
	if err != nil {
		return err
	}
	return nil
}
