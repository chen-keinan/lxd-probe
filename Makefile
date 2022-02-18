SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv ~/vms/lxd-probelxd-probe ~/vms-local/lxd-probe
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=lxd-probe
GOCOPY=cp lxd-probe ~/vagrant_file/.

all:test lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md
build:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/lxd-probe;
build_local:
	export PATH=$GOPATH/bin:$PATH;
	export PATH=$PATH:/home/vagrant/go/bin
	export PATH=$PATH:/home/root/go/bin
	$(GOBUILD) ./cmd/lxd-probe;
install:build_travis
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
test_build_travis:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) -short ./...  -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/lxd-probe;
build_travis:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/lxd-probe;
build_remote:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/lxd-probe
	mv lxd-probe ~/boxes/basic_box/lxd-probe

build_docker_local:
	docker build -t chenkeinan/lxd-probe:3 .
	docker push chenkeinan/lxd-probe:3
dlv:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./lxd-probe
build_beb:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/lxd-probe/lxd-probe.go
	scripts/deb.sh
.PHONY: all build install test
