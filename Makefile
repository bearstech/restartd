.PHONY: restartctl restartd gopath

export GOPATH:=$(shell pwd)/gopath

all: restartctl restartd

gopath/src/github.com/bearstech/restartd:
	mkdir -p gopath/src/github.com/bearstech
	ln -s ../../../.. gopath/src/github.com/bearstech/restartd

gopath/src/gopkg.in/yaml.v2:
	go get gopkg.in/yaml.v2

gopath/src/github.com/Sirupsen/logrus:
	go get github.com/Sirupsen/logrus

gopath/src/github.com/golang/protobuf/proto:
	go get github.com/golang/protobuf/proto

gopath: gopath/src/github.com/bearstech/restartd

deps: gopath/src/gopkg.in/yaml.v2 gopath/src/github.com/Sirupsen/logrus gopath/src/github.com/golang/protobuf/proto

bin:
	mkdir -p bin

restartctl: bin gopath deps
	go build -o bin/restartctl github.com/bearstech/restartd/restartctl/

restartd: bin gopath deps
	go build -o bin/restartd github.com/bearstech/restartd/restartd/

test:
	go test github.com/bearstech/restartd/restartctl/
	go test github.com/bearstech/restartd/protocol/
	go test github.com/bearstech/restartd/restartd/

install:
	cp bin/restartd /usr/local/bin
	cp bin/restartctl $(ROOTFS)/usr/local/bin

clean:
	rm -rf gopath
	rm -rf bin

linux:
	docker run -it --rm -v `pwd`:/go golang make

protoc:
	protoc --go_out=. protocol/*.proto

vet:
	go vet github.com/bearstech/restartd/restartctl
	go vet github.com/bearstech/restartd/restartd
