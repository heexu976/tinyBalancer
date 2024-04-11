SERVER_VERSION = v1.0.0

GO           = go
GO_VENDOR    = go mod
MKDIR_P      = mkdir -p

.PHONY: build
build:
	GO111MODULE=on $(GO) build -v -o _output/tinyBalancer ./

.PHONY: docker.build
docker.build: 
	docker build --no-cache --rm --tag tinybalancer:$(SERVER_VERSION) -f ./build/Dockerfile .