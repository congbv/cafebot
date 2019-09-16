run:
	go run ./*.go --log-level=debug

test:
	go test --race ./... 
