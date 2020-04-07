.DEFAULT_GOAL := all

.PHONY: all test tests bin build clean clean_release clean_release

all: test
test:
	@go test -gcflags=-l ./...

tests: test

build: bin

clean:
	@rm -r bin

clean_release:
	@rm -rf release/

release: test clean_release
	@mkdir -p release/udp
	@go build -ldflags="-s -w" -o release/udp/bin
	@echo '"release" successful'


bin: test
	@mkdir -p bin
	@go build -o bin/udp-wejay
	@echo "\"bin/udp-wejay\" created"

