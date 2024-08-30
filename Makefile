GOPATH := $(shell go env GOPATH)

install:
	go mod tidy
run:
	go run main.go start
check:
	go run main.go check -d=2

.PHONY: install run check
