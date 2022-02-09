# ledger


protoc ./api/*.proto --go_out=api --go-grpc_out=api

docker-compose build --no-cache

go test -p 1 .\e2e\... -v