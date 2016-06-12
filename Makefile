.PHONY: restartctl restartd

all: restartctl restartd

restartctl:
	cd restartctl && go build
restartd:
	cd restartd && go build
