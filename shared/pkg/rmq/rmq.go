package rmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}

func (r *RabbitMQ) SetPrefetchCount(num int) error {
	return r.ch.Qos(num, 0, false)
}

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) ConsumeQueue(queue amqp.Queue) (<-chan amqp.Delivery, error) {
	return r.ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Publish(
	ctx context.Context,
	queue amqp.Queue,
	body []byte,
	correlationId *string,
	replyTo *string,
) error {
	return r.ch.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: *correlationId,
			ReplyTo:       *replyTo,
		},
	)
}
