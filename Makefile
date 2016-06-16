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

clean:
	rm -f restartctl/restartctl
	rm -f restartd/restartd

linux:
	docker run -it --rm -v `pwd`:/go golang make
