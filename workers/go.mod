module app.workers

go 1.24.2

require (
	app.shared v0.0.0
	github.com/fatih/color v1.18.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace app.shared => ../shared
