.PHONY: default
default: bin/hw bin/hwc bin/hwl

VERSION := $(shell git describe --tags 2> /dev/null || echo "unreleased")
V_DIRTY := $(shell git describe --exact-match HEAD 2> /dev/null > /dev/null || echo "-unreleased")
bin/hw: cmd/hw/main.go bin
	@GOOS=linux go build -ldflags \
		"-X mcquay.me/hw.Version=$(VERSION)$(V_DIRTY)" \
		 -v -o bin/hw ./cmd/hw

bin/hwc: cmd/hwc/main.go bin
	@GOOS=linux go build -ldflags \
		"-X mcquay.me/hw.Version=$(VERSION)$(V_DIRTY)" \
		 -v -o bin/hwc ./cmd/hwc

bin/hwl: cmd/hwl/main.go bin
	@GOOS=linux go build -ldflags \
		"-X mcquay.me/hw.Version=$(VERSION)$(V_DIRTY)" \
		 -v -o bin/hwl ./cmd/hwl

bin:
	mkdir bin
