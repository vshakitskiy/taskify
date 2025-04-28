build:
	cd server && go build -o ../bin/server ./cmd/main.go && cd .. && \
	cd client && go build -o ../bin/client ./cmd/main.go && cd .. && \
	cd workers && go build -o ../bin/worker ./cmd/main.go && cd ..