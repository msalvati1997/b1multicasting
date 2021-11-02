#  Makefile for Multicast project

all:    ## clean, format, build and unit test
	make clean
	make build
	make test
clean:  ## go clean
	go clean
build:
	docker-compose up
test:
	go test -v ./pkg/multicasting/test
test-api:
	go test -v ./pkg/multicastapp/testingapi
