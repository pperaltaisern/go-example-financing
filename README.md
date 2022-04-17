# go-example-financing

This is an example on how to build event sourced and CQRS applications in go.
- Event sourcing library: esrc (https://github.com/pperaltaisern/esrc)
- CQRS library: watermill (https://github.com/ThreeDotsLabs/watermill)

## Architecture
![](doc/architecture.jpg)

## Acceptance criteria

## Demo

![](doc/demo.gif)

## DB state

## Code

## Testing strategy

## Considerations

## Dev

Run:
```bash
docker-compose up 
go run cmd/populate/main.go # publishes integration events for registering issuers and investors
```
Use a grpc client like [bloomrpc](https://github.com/bloomrpc/bloomrpc) to interact with the application. Send commands to port :8080 and queries to :8081.


Unit test:
```bash
go test ./... -short
```

Acceptance test:
```bash
docker-compose up -d
go test ./acceptance/... -v -p 1 -count 1
```

Gen proto:
```bash
protoc ./api/*.proto --go_out=api --go-grpc_out=api
```
