@echo off
docker build ./ -t mcs/gomicro-service:V1
docker run --name gomicro-service -p 9543:8443 -p 9080:8080 mcs/gomicro-service:V1