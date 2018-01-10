bin/hw: main.go bin
	GOOS=linux go build -v -o bin/hw

bin:
	mkdir bin
