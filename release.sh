#!/bin/bash

docker-compose build
docker tag sdpequinor/release-aware-web:latest sdpequinor/release-aware-web:$1
docker push sdpequinor/release-aware-web:$1

docker tag sdpequinor/release-aware-api:latest sdpequinor/release-aware-api:$1
docker push sdpequinor/release-aware-api:$1