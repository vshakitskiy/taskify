proto_gen:
	protoc --proto_path api/proto \
	--go_out=pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	api/proto/tasks.proto