.PHONY: test

test:
	go test -v ./tests

build:
	go build -o API

run: build
	go build -o API
	./API