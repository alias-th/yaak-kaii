PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/**/*.proto)
GO_OUT := .

generate_proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

random_key:
	openssl rand -base64 32

.PHONY: generate_proto random_key