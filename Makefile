.PHONY: restartctl restartd

all: restartctl restartd

restartctl:
	cd restartctl && go build
restartd:
	cd restartd && go build

test:
	cd protocol && go test
	cd restartd && go test

get:
	go get gopkg.in/yaml.v2

install:
	cp restartd/restartd /usr/local/bin
	cp restartctl/restartctl $(ROOTFS)/usr/local/bin

clean:
	rm -f restartctl/restartctl
	rm -f restartd/restartd

linux: src/gopkg.in/yaml.v2
	docker run -it --rm -v $(GOPATH):/go -w /go/src/restartd golang make

src/gopkg.in/yaml.v2:
	docker run -it --rm -v $(GOPATH):/go -w /go/src/restartd golang make get
