.PHONY: default
default: bin/hw bin/hwc bin/hwl

bin/hw: cmd/hw/main.go bin
	GOOS=linux go build -v -o bin/hw ./cmd/hw

bin/hwc: cmd/hwc/main.go bin
	GOOS=linux go build -v -o bin/hwc ./cmd/hwc

bin/hwl: cmd/hwl/main.go bin
	GOOS=linux go build -v -o bin/hwl ./cmd/hwl

bin:
	mkdir bin
