
# Makefile for github-access-checker

PROTO_DIR=proto/
OUT_DIR=proto/
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

all: generate

.PHONY: install
install:
	go install tool

.PHONY: generate
generate:
	@echo "Generating protobuf code..."
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt paths=source_relative \
		--openapiv2_out ./openapiv2 --openapiv2_opt logtostderr=true \
		$(PROTO_FILES)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run:
	go run main.go

.PHONY: build
run:
	go build

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)
