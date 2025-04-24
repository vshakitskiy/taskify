module app.workers

go 1.24.2

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	shared v0.0.0
)

replace shared => ../shared
