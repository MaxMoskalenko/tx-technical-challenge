IMAGE_NAME ?= tx-technical-challenge
GRPC_PORT ?= 50051

.PHONY: proto docker-build docker-run

proto:
	protoc -I . internal/proto/api.proto --go_out=. --go_opt=module=github.com/MaxMoskalenko/tx-technical-challenge  --go-grpc_out=. --go-grpc_opt=module=github.com/MaxMoskalenko/tx-technical-challenge

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run --rm -p $(GRPC_PORT):50051 -e GRPC_ADDR=:50051 $(IMAGE_NAME)
