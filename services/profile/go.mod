module profile

go 1.26

require (
	github.com/lib/pq v1.12.3
	github.com/segmentio/kafka-go v0.4.51
	jwt v0.0.0
	middlewares v0.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace (
	jwt => ../../pkg/jwt
	logger => ../../pkg/logger
	middlewares => ../../pkg/middlewares
)
