.DEFAULT_GOAL := all

.PHONY: all test bin build clean

all: test
test:
	@go test ./...

build: bin

clean:
	@rm -r bin

bin: test
	@mkdir -p bin
	@go build -o bin/udp-wejay
	@echo "\"bin/udp-wejay\" created"

