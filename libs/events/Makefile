PROTO_PATH := .
PROTOFILES := $(shell find $(PROTO_PATH) -type f -name '*.proto')

.PHONY: check-clang-format
check-clang-format:
	@type clang-format >/dev/null 2>&1 || (echo "clang-format has not been found, please install it from brew or your package manager" && false)

.PHONY: format
format: check-clang-format
	@echo Reformatting protobuf files
	@clang-format -style=file -i ${PROTOFILES}

.PHONY: build
build:
	@protoc --proto_path=. --go_out=. --go_opt=paths=source_relative ${PROTOFILES}
