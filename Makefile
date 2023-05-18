all: format test vet lint generate build bin

generate:
	go generate ./...

format:
	go fmt ./...

lint:
	${GOPATH}/bin/staticcheck ./...

vet:
	go vet ./...

test:
	go test ./...

build: generate format lint vet test
	go build ./...

bin: generate
	-go build cmd/stampWalletServer.go
