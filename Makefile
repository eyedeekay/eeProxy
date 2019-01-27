
GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

build:
	go build $(GO_COMPILER_OPTS)

all: linux osx windows

linux:
	GOOS=linux go build $(GO_COMPILER_OPTS) -o eeProxy

osx:
	GOOS=darwin go build $(GO_COMPILER_OPTS) -o eeProxy.app

windows:
	GOOS=windows go build $(GO_COMPILER_OPTS) -o eeProxy.exe

fmt:
	find . -name '*.go' -exec gofmt -w {} \;

clean:
	go clean

deps:
	go get -u github.com/eyedeekay/eeproxy
	go get -u github.com/eyedeekay/go-socks5

test:
	mkdir -p testdir && cd testdir && \
		../eeProxy

kill:
	killall -9 eeProxy

fire:
	\/usr/bin/curl --socks5-hostname 127.0.0.1:7950 http://i2p-projekt.i2p --output i2p-projekt.html
	\/usr/bin/curl --socks5-hostname 127.0.0.1:7950 http://inr.i2p --output inr.html

tb:
	go build ./conn
	go build ./socks
	go build ./rewriter
	go build ./resolve

firefox:
	rm -rf ./testdir/profile.i2p && cp -rv ./profile.i2p ./testdir/profile.i2p
	firefox --profile ./testdir/profile.i2p http://inr.i2p
