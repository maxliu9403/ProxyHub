PROJECT=gin-template

GOCMD=CGO_ENABLED=0 go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=$(PROJECT)-api
BINARY_WIN=$(BINARY_NAME).exe

VERSION=$(shell git describe --abbrev=0 --tags)
COMPILER=$(shell go version)
DATE=$(shell date)
COMMITID=$(shell git log --pretty=format:"%h" -1)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

default: build

build:
	mkdir -p bin/
	$(GOBUILD) -o bin/$(BINARY_NAME) -ldflags " \
			-extldflags '-static' \
			-X 'main._version=$(VERSION)' \
			-X 'main._date=$(DATE)' \
			-X 'main._commit=$(COMMITID)' \
			-X 'main._compiler=$(COMPILER)'" cmd/app/main.go

build-vendor:
	mkdir -p bin/
	$(GOBUILD) -mod=vendor -o bin/$(BINARY_NAME) -ldflags " \
			-extldflags '-static' \
			-X 'main._version=$(VERSION)' \
			-X 'main._date=$(DATE)' \
			-X 'main._commit=$(COMMITID)' \
			-X 'main._compiler=$(COMPILER)'" cmd/app/main.go

test:
	$(GOTEST)  -v -run=. ./...

clean:
	echo "clean"

install:
	echo "install"

doc:
	action-swag init -g ./cmd/app/main.go

.PHONY: default build clean install test