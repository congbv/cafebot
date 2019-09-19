run:  vet
	go run ./*.go --log-level=debug

vet:
	go vet ./...

test:
	go test --race ./... 
