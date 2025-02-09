.PHONY: build clean test docker run

GO=CGO_ENABLED=0 GOOS=linux go
GOCGO=CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go

cmd/hummingbird-core/hummingbird-core:
	$(GOCGO) build -ldflags "-s -w" -o $@ ./cmd/hummingbird-core

cmd/mqtt-broker/mqtt-broker:
	$(GO) build -ldflags "-s -w" -o $@ ./cmd/mqtt-broker

generate/api:
	cd cmd/hummingbird-core && swag init --parseDependency --parseInternal --parseDepth 10

.PHONY: start
start:
	go run cmd/hummingbird-core/main.go -c cmd/hummingbird-core/res/configuration.toml

.PHONY: build
build:
	docker buildx build --platform linux/amd64 -t "registry.cn-shanghai.aliyuncs.com/winc-link/hummingbird-core:v1.0" -f cmd/hummingbird-core/Dockerfile --push .

