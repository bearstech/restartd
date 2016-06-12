.PHONY: restartctl restartd

all: restartctl restartd

restartctl:
	cd restartctl && go build
restartd:
	cd restartd && go build

test:
	cd protocol && go test
	cd restartd && go test

clean:
	rm -f restartctl/restartctl
	rm -f restartd/restartd
