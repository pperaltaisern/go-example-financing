# ledger


protoc ./api/*.proto --go_out=api --go-grpc_out=api

docker-compose build --no-cache