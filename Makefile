V=`git rev-parse --short HEAD`
B="-X main.Version $(V)"

run:
	@godep go run main.go

test:
	@godep go test -p 1 -cover -race -bench=. -benchmem ./...

build:
	@godep go build -ldflags=$(B) -o bin/forward-ssh

.PHONY: test build
