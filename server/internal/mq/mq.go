package mq

import (
	"fmt"

	"app.shared/pkg/rmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Methods      *rmq.RabbitMQ
	TaskQueue    amqp.Queue
	ResultsQueue amqp.Queue
	ResultsCh    <-chan amqp.Delivery
}

func NewMQ(url string, prefetchCount int) (*RabbitMQ, error) {
	mq, err := rmq.NewRabbitMQ(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %s", err)
	}

	resultsQueue, err := mq.DeclareQueue("results")
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %s", err)
	}

	err = mq.SetPrefetchCount(prefetchCount)
	if err != nil {
		return nil, fmt.Errorf("Failed to set QoS: %s", err)
	}

	tasksChannel, err := mq.ConsumeQueue(resultsQueue)
	if err != nil {
		return nil, fmt.Errorf("Failed to consume a queue: %s", err)
	}

	tasksQueue, err := mq.DeclareQueue("tasks")
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %s", err)
	}

	return &RabbitMQ{
		Methods:      mq,
		TaskQueue:    tasksQueue,
		ResultsQueue: resultsQueue,
		ResultsCh:    tasksChannel,
	}, nil
}
