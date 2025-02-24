run : build
	@./bin/go-redis --listenAddr :8080

build:
	@go build -o bin/go-redis
