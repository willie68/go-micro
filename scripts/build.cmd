@echo off
go build -ldflags="-s -w" -o gomicro-service.exe cmd/service/main.go