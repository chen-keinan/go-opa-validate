SHELL := /bin/bash

GOCMD=go
GOMOCKS=$(GOCMD) generate ./...
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test

all:
	$(info  "completed running make file for go-opa-validate")
fmt:
	@go fmt ./...
lint:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen
	export PATH=$HOME/go/bin:$PATH
	$(GOMOCKS)
	lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md

.PHONY: install-req fmt lint tidy test imports .
