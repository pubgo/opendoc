vet:
	go vet ./...

run-example:
	go run internal/examples/*.go

goimports:
	goimports -w -local github.com/pubgo/opendoc .

refactor:
	gofumpt -l -w -extra .

lint:
	golangci-lint run --timeout=10m --verbose
