@echo off
docker build ./ -t easy_software/gomicro-service-go:V1
docker run --name gomicro-service-go -p 9443:9443 -p 9080:9080 easy_software/gomicro-service-go:V1