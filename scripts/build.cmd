@echo off
swag init -g ./cmd/service/main.go
swag init -d "./internal/apiv1" -g "../../cmd/service/main.go" --parseDependency --parseDepth 2 -o "./api"
go build -ldflags="-s -w" -o gomicro-service.exe cmd/service/main.go