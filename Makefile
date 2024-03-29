VERSION := $(shell git describe --tags)

BINPATH := ./bin

GOCMD := go
GOBUILD := $(GOCMD) build
GOINSTALL := $(GOCMD) install
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
BINARY_NAME := luks2crypt
MOCKSERVER_NAME := cryptservermock

VAGRANTCMD := vagrant

LDFLAGS=-ldflags "-X main.VERSION=$(VERSION)"

all: test build

install:
	$(GOINSTALL) $(LDFLAGS) -v ./cmd/$(BINARY_NAME)

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINPATH)/$(BINARY_NAME) -v ./cmd/$(BINARY_NAME)

lint:
	go vet ./...

test:
	$(GOTEST) -race -v ./...

clean:
	$(GOCLEAN)
	rm -r ./bin ./tmp ./artifacts ./dist

deps:
	$(GOMOD) tidy
	$(GOCMD) get -u ./...

build-mockserver:
	$(GOCMD) build -o $(BINPATH)/$(MOCKSERVER_NAME) -v ./tools/cryptservermock

mockserver: build-mockserver
	sudo $(BINPATH)/$(MOCKSERVER_NAME)

devup:
	$(VAGRANTCMD) up

devssh:
	$(VAGRANTCMD) ssh

devclean:
	$(VAGRANTCMD) destroy --force
