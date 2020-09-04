
ARCH ?= $(shell go env GOARCH)
VERSION ?= test

.PHONY: build
build:
	GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 go build -o _output/$(ARCH)/pod-exec cmd/web.go

.PHONY: image
image:
	docker build --network host -t pod-exec:$(VERSION) -f build/Dockerfile .
	docker tag pod-exec:$(VERSION) hub-cn-shanghai-2.kce.ksyun.com/ksyun/pod-exec:$(VERSION)
	docker push hub-cn-shanghai-2.kce.ksyun.com/ksyun/pod-exec:$(VERSION)