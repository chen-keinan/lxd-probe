SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv ~/vms/lxd-probelxd-probe ~/vms-local/lxd-probe
GOPACKR=$(GOCMD) get -u github.com/gobuffalo/packr/packr && packr
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
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md
build:
	$(GOPACKR)
	export PATH=$GOPATH/bin:$PATH;
	export PATH=$PATH:/home/vagrant/go/bin
	export PATH=$PATH:/home/root/go/bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/lxd-probe;
install:build_travis
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
test_travis:
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
build_travis:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/lxd-probe;
build_remote:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' ./cmd/lxd-probe
	mv lxd-probe ~/boxes/box_ebpf//lxd-probe
build_docker:
	export PATH=$GOPATH/bin:$PATH;
	docker build -t chenkeinan/lxd-probe:latest .
	docker push chenkeinan/lxd-probe:latest
build_docker_local:
	docker build -t chenkeinan/lxd-probe:latest .
	docker push chenkeinan/lxd-probe:latest
dlv:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./lxd-probe
build_beb:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/lxd/lxd-probe.go
	scripts/deb.sh
.PHONY: all build install test
