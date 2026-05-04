BIN_DIR := $(CURDIR)/bin

.DEFAULT_GOAL := run

.PHONY: build run test controller worker cli watch

controller:
	go run $(CURDIR)/cmd/controller

worker:
	go run $(CURDIR)/cmd/worker

cli:
	go run $(CURDIR)/cmd/cli

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/controller $(CURDIR)/cmd/controller
	go build -o $(BIN_DIR)/worker $(CURDIR)/cmd/worker
	go build -o $(BIN_DIR)/cli $(CURDIR)/cmd/cli

watch:
	air

run:
	@make -j 3 controller worker cli