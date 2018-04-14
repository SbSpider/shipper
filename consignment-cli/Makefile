build:
	GOOS=linux GOARCH=amd64 go build
	docker build -t consignment-cli .

run: build
	docker run -e MICRO_REGISTRY=mdns consignment-cli
