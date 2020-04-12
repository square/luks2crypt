VERSION := $(shell git describe --tags)

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell dpkg --print-architecture)

GITTAG := $(shell git describe --tags)

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

deploytar:
	mkdir -p tmp/$(BINARY_NAME) artifacts
	cp bin/$(BINARY_NAME) README.md COPYING LICENSE.txt tmp/$(BINARY_NAME)/
	tar -C tmp -czvf artifacts/$(BINARY_NAME)-$(GITTAG)-$(OS)-$(ARCH).tar.gz $(BINARY_NAME)

deploy:
	GO111MODULE=off go get github.com/tcnksm/ghr
	ghr -t ${GITHUB_TOKEN} -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} -c ${CIRCLE_SHA1} -delete ${CIRCLE_TAG} ./artifacts/

lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	golint -set_exit_status ./...
	go vet ./...

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -r ./bin ./tmp ./artifacts

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
