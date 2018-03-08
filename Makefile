bin/hw: cmd/hw/main.go bin
	GOOS=linux go build -v -o bin/hw ./cmd/hw

bin:
	mkdir bin
