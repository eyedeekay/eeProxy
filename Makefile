
build:
	go build

fmt:
	find . -name '*.go' -exec gofmt -w {} \;

clean:
	go clean
