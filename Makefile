# This Makefile is basically for unit testing only
# To build and install this library, please use normal commands
# such as "go get" or "go install"

export GOPATH=$(shell pwd)/_gopath
export EXAMPLE=$(shell pwd)/_examples

all: fmt test

fmt:
	@echo "Format the source files"
	@echo "-----------------------"
	go fmt
	cd _examples && go fmt
	@echo

test: test.unit test.example

test.example: gofan _examples/run-all
	@echo "Run Examples"
	@echo "------------"
	@_examples/run-all
	@echo

test.unit: gofan
	@echo "Unit Tests"
	@echo "----------"
	go version
	go test
	@echo

clean:
	rm -Rf _gopath/*

gofan: _gopath/src/github.com/yookoala/gofan
	@echo "Install gofan"
	@echo "-------------"
	rm -Rf _gopath/pkg/*/github.com/yookoala
	go install github.com/yookoala/gofan
	@echo

_examples/run-all: \
	gofan
	@echo "Build Example(s) runner"
	@echo "-----------------------"
	cd _examples && go build -o ${EXAMPLE}/run-all
	@echo

_gopath/src:
	@echo "Create testing GOPATH"
	@echo "---------------------"
	mkdir -p _gopath/src
	@echo

_gopath/src/github.com/yookoala/gofan:
	@mkdir -p _gopath/src/github.com/yookoala
	@cd _gopath/src/github.com/yookoala && ln -s ../../../../. gofan

.PHONY: test test.unit test.example gofan clean
