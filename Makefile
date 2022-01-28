SHELL := /bin/bash

GOCMD=go
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test

all:
	$(info  "completed running make file for go-opa-validate")
fmt:
	@go fmt ./...
lint:
	./lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) install github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen
	export PATH=$GOPATH/bin:$PATH;
	$(GOCMD) generate ./...
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md

.PHONY: install-req fmt lint tidy test imports .
