@echo off
swag init -g ./cmd/service/main.go
go build -ldflags="-s -w" -o gomicro-service.exe cmd/service/main.go