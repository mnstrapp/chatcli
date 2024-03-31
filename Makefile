SHELL = /bin/bash

PREFIX_PATH := ${HOME}/.local
TEMP_PATH := /tmp/mnstrapp/chatcli
CWD_PATH = $(shell pwd)

.PHONY: all
all: clean deps gen build

deps: deps_structure deps_grpc

deps_structure:
	mkdir -p ${PREFIX_PATH}
	mkdir -p ${TEMP_PATH}

deps_grpc:
	if ! command -v protoc &> /dev/null; then \
		if [ ! -d "${TEMP_PATH}/grpc" ]; then \
			git clone \
				--recurse-submodules \
				-b v1.62.0 \
				--depth 1 \
				--shallow-submodules https://github.com/grpc/grpc \
				${TEMP_PATH}/grpc; \
		fi; \
		cd ${TEMP_PATH}/grpc/ && \
			mkdir -p ${TEMP_PATH}/grpc/cmake/build/ && \
			pushd ${TEMP_PATH}/grpc/cmake/build/ && \
			cmake -DgRPC_INSTALL=ON \
			-DgRPC_BUILD_TESTS=OFF \
			-DCMAKE_INSTALL_PREFIX=${PREFIX_PATH} \
			../.. && \
			make -j4 && \
			make install; \
	fi
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

gen:
	protoc \
		--proto_path=${CWD_PATH}/proto \
		--go_out=${CWD_PATH}/gossip \
		--go-grpc_out=${CWD_PATH}/gossip \
		--plugin=protoc-gen-go=$(shell go env GOPATH)/bin/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(shell go env GOPATH)/bin/protoc-gen-go-grpc \
		gossip.proto

build:
	go build -o ${CWD_PATH}/target/chatcli

install: build
	cp target/chatcli ${PREFIX_PATH}/bin/

uninstall:
	if [ -f "${PREFIX_PATH}/bin/chatcli" ]; then \
		rm -rf ${PREFIX_PATH}/bin/chatcli; \
	fi

clean:
	rm -rf ${CWD_PATH}/target
	rm -rf ${TEMP_PATH}