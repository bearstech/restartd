.PHONY: restartctl restartd gopath

export GOPATH:=$(shell pwd)/gopath
GOTOOLS = \
	golang.org/x/tools/cmd/cover \
	github.com/axw/gocov/gocov \
	gopkg.in/matm/v1/gocov-html

GITCOMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS := ${LDFLAGS} \
	-X main.GITCOMMIT=${GITCOMMIT}

all: restartctl restartd

gopath/src/github.com/bearstech/restartd:
	mkdir -p gopath/src/github.com/bearstech
	ln -s ../../../.. gopath/src/github.com/bearstech/restartd

gopath/src/gopkg.in/yaml.v2:
	go get gopkg.in/yaml.v2

gopath/src/github.com/Sirupsen/logrus:
	go get github.com/Sirupsen/logrus

gopath/src/github.com/coreos/go-systemd/dbus:
	go get github.com/coreos/go-systemd/dbus

gopath/src/github.com/golang/protobuf/proto:
	go get github.com/golang/protobuf/proto

gopath/src/github.com/urfave/cli:
	go get github.com/urfave/cli

gopath: gopath/src/github.com/bearstech/restartd

deps: gopath/src/gopkg.in/yaml.v2 gopath/src/github.com/Sirupsen/logrus gopath/src/github.com/golang/protobuf/proto gopath/src/github.com/coreos/go-systemd/dbus gopath/src/github.com/urfave/cli

bin:
	mkdir -p bin

restartctl: bin gopath deps
	go build -ldflags "${LDFLAGS}" -o bin/restartctl github.com/bearstech/restartd/cli/restartctl/

restartd: bin gopath deps
	go build -ldflags "${LDFLAGS}" -o bin/restartd github.com/bearstech/restartd/cli/restartd/

test: gopath/src/github.com/bearstech/restartd deps
	go test github.com/bearstech/restartd/listen/
	go test github.com/bearstech/restartd/protocol/
	go test github.com/bearstech/restartd/restartd/
	go test github.com/bearstech/restartd/model/

install:
	cp bin/restartd /usr/local/sbin
	cp bin/restartctl $(ROOTFS)/opt/factory/

clean:
	rm -rf gopath
	rm -rf bin

linux:
	docker run -it --rm -v `pwd`:/go golang make

protoc:
	protoc --go_out=. model/*.proto

vet:
	go vet github.com/bearstech/restartd/restartctl
	go vet github.com/bearstech/restartd/restartd

tools:
	go get -u -v $(GOTOOLS)

cover:
	go test github.com/bearstech/restartd/listen/ --cover
	go test github.com/bearstech/restartd/restartctl/ --cover
	go test github.com/bearstech/restartd/protocol/ --cover
	go test github.com/bearstech/restartd/restartd/ --cover
