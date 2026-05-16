BIN_DIR := "bin"
PROTO_DIR := "proto"
OUT_DIR := "pkg/protocol/v1"

run:
	@just -j 2 controller worker

controller:
	go run ./cmd/controller

worker:
	go run ./cmd/worker

cli:
	go run ./cmd/cli

build:
	@mkdir -p {{BIN_DIR}}
	go build -o {{BIN_DIR}}/controller ./cmd/controller
	go build -o {{BIN_DIR}}/worker ./cmd/worker
	go build -o {{BIN_DIR}}/cli ./cmd/cli

watch:
	air

proto:
	protoc --go_out={{OUT_DIR}} --go_opt=paths=source_relative \
	--go-grpc_out={{OUT_DIR}} --go-grpc_opt=paths=source_relative \
	-I/usr/include \
	--proto_path={{PROTO_DIR}} \
	{{PROTO_DIR}}/orchestrator.proto


image_build:
	docker build -t codebuddha/orchestrator:latest .
image_push:
	docker push codebuddha/orchestrator:latest


start:
	docker compose -f docker-compose.dev.yml up --build