run: lint vet
	go run ./*.go --log-level=debug

lint:
	golint ./... 

vet:
	go vet ./...

test:
	go test --race ./... 
