# ledger


protoc ./api/*.proto --go_out=api --go-grpc_out=api

docker-compose build --no-cache

go test ./acceptance/... -v -run TestQueryFeatures -p 1 -count 1
go test ./acceptance/... -v -run TestCommandFeatures -p 1 -count 1
go test ./acceptance/... -v -p 1 -count 1