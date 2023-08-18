package mq

import (
	"github.com/streadway/amqp"
	"log"
)

type ReplenishmentConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewReplenishmentConsumer(url string) (*ReplenishmentConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &ReplenishmentConsumer{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *ReplenishmentConsumer) Consume() error {
	queue, err := c.channel.QueueDeclare(
		"replenishment", // queue name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	messages, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
			// Call your replenishment balance logic here using d.Body
		}
	}()

	log.Printf("Waiting for messages. To exit, press Ctrl+C")
	<-forever

	return nil
}
