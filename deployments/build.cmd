@echo off
go build -ldflags="-s -w" -o serice-gomicro-go.exe cmd/service.go