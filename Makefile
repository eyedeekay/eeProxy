
build:
	go build

fmt:
	find . -name '*.go' -exec gofmt -w {} \;

clean:
	go clean

deps:
	go get -u github.com/eyedeekay/eeproxy
