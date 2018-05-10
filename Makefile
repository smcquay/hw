.PHONY: default
default: bin/hw bin/hwc bin/hwl

VERSION := $(shell git describe --tags 2> /dev/null || echo "unreleased")
V_DIRTY := $(shell git describe --exact-match HEAD 2> /dev/null > /dev/null || echo "-unreleased")
GIT     := $(shell git rev-parse --short HEAD)
DIRTY   := $(shell git diff-index --quiet HEAD 2> /dev/null > /dev/null || echo "-dirty")

bin/hw: $(shell ls *.go) cmd/hw/main.go bin
	@echo hw
	@GOOS=linux go build -ldflags \
		"-X mcquay.me/hw.Version=$(VERSION)$(V_DIRTY) \
		 -X mcquay.me/hw.Git=$(GIT)$(DIRTY)" \
		 -v -o bin/hw ./cmd/hw

bin/hwc: $(shell ls *.go ) cmd/hwc/main.go bin
	@echo hwc
	@GOOS=linux go build -ldflags \
		"-X mcquay.me/hw.Version=$(VERSION)$(V_DIRTY) \
		 -X mcquay.me/hw.Git=$(GIT)$(DIRTY)" \
		 -v -o bin/hwc ./cmd/hwc

bin/hwl: $(shell ls *.go) cmd/hwl/main.go bin
	@echo hwl
	@GOOS=linux go build -ldflags \
		"-X mcquay.me/hw.Version=$(VERSION)$(V_DIRTY) \
		 -X mcquay.me/hw.Git=$(GIT)$(DIRTY)" \
		 -v -o bin/hwl ./cmd/hwl

bin:
	mkdir bin

.PHONY: clean
clean:
	@rm -fv bin/{hw,hwl,hwc}

.PHONY: lint
lint:
	@golint $(shell go list mcquay.me/hw/...)
	@go vet $(shell go list mcquay.me/hw/...)

.PHONY: docker-build
docker-build: bin/hw bin/hwc bin/hwl
	docker build -f cmd/hw/Dockerfile  . -t smcquay/hw:$(VERSION)$(V_DIRTY)
	docker build -f cmd/hwc/Dockerfile . -t smcquay/hwc:$(VERSION)$(V_DIRTY)
	docker build -f cmd/hwl/Dockerfile . -t smcquay/hwl:$(VERSION)$(V_DIRTY)

.PHONY: docker-push
docker-push:
	docker push smcquay/hw:$(VERSION)$(V_DIRTY)
	docker push smcquay/hwc:$(VERSION)$(V_DIRTY)
	docker push smcquay/hwl:$(VERSION)$(V_DIRTY)
