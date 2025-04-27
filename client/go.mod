module app.client

go 1.24.2

require (
	github.com/fatih/color v1.18.0
	google.golang.org/grpc v1.72.0
	app.server v0.0.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250425173222-7b384671a197 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace app.server => ../server