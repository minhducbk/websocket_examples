GO=$(shell which go)
WORK_DIR=.
.DEFAULT_GOAL := build

.PHONY: test
## Run all unit tests
test:
	${GO} test `${GO} list ./...` -coverprofile coverage.out

.PHONY: build
## Build all binaries
build:
	mkdir -p bin/
	${GO} build -o bin/chat chat/main.go 
	${GO} build -o bin/ws_deribit ws_deribit/main.go 

.PHONY: chat
## Run chat
chat:
	${GO} run chat/*.go 

.PHONY: ws_coinbase
## Run chat
ws_coinbase:
	${GO} run ws_coinbase/*.go 

.PHONY: ws_coincap
## Run chat
ws_coincap:
	${GO} run ws_coincap/*.go 

.PHONY: ws
## Run chat
ws:
	${GO} run ws/*.go 
