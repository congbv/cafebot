install: vet
	go install -ldflags="-w -s" ./...

run:  vet test
	go run ./*.go --log-level=debug
.PHONY: run

vet:
	go vet ./...
.PHONY: vet

test:
	go test --race ./... 
.PHONY: test


