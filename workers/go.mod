module workers

go 1.24.2

require shared v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
)

replace shared => ../shared
