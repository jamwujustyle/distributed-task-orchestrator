BIN_DIR := $(CURDIR)/bin
PROTO_DIR := $(CURDIR)/proto
OUT_DIR := $(CURDIR)/pkg/protocol/v1

.DEFAULT_GOAL := run

.PHONY: build run test controller worker cli watch proto

controller:
	go run "$(CURDIR)/cmd/controller"

worker:
	go run "$(CURDIR)/cmd/worker"

cli:
	go run "$(CURDIR)/cmd/cli"

build:
	@mkdir -p "$(BIN_DIR)"
	go build -o "$(BIN_DIR)/controller" "$(CURDIR)/cmd/controller"
	go build -o "$(BIN_DIR)/worker" "$(CURDIR)/cmd/worker"
	go build -o "$(BIN_DIR)/cli" "$(CURDIR)/cmd/cli"

watch:
	air

run:
	@make -j 3 controller worker cli

proto:
	protoc --go_out="$(OUT_DIR)" --go_opt=paths=source_relative \
	--go-grpc_out="$(OUT_DIR)" --go-grpc_opt=paths=source_relative \
	-I/usr/include \
	--proto_path="$(PROTO_DIR)" \
	"$(PROTO_DIR)/orchestrator.proto"
