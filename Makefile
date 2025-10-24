vet:
	go vet ./...

run-example:
	go run internal/examples/*.go

goimports:
	goimports -w -local github.com/pubgo/opendoc .
