embed:
	embedmd -w README.md

example:
	go run .docs/examples/quickstart.go

test:
	go test ./...