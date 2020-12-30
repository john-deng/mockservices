#!/usr/bin/env bash

if [[ $(docker network ls | grep -c solarmesh) == 0 ]]; then
  docker network create solarmesh
fi


function clean_up_docker() {
  container=$1
  if [[ $(docker ps | grep -c "${container}" ) != 0 ]]; then
    docker rm -f "${container}"
  fi
}

clean_up_docker product
# Run as product v1
docker run -it -d \
  --name=product \
  --net=solarmesh \
  -p 8080:8080 \
  -e APP_NAME=product \
  -e APP_VERSION=v1 \
  -e CLUSTER_NAME=cluster01 \
  solarmesh/solar-mock-app:latest

clean_up_docker inventory
# Run as inventory v1
docker run -it -d \
  --name=inventory \
  --net=solarmesh \
  -p 8081:8080 \
  -e APP_NAME=inventory \
  -e APP_VERSION=v1 \
  -e CLUSTER_NAME=cluster01 \
  -e UPSTREAM_URLS='http://product:8080/,' \
  solarmesh/solar-mock-app:latest

clean_up_docker payment
# Run as payment v2
docker run -it -d \
  --name=payment \
  --net=solarmesh \
  -p 8082:8080 \
  -e APP_NAME=payment \
  -e APP_VERSION=v2 \
  -e CLUSTER_NAME=cluster01 \
  -e UPSTREAM_URLS='http://product:8080/,http://inventory:8080/,' \
  solarmesh/solar-mock-app:latest

clean_up_docker order
# Run as order v1
docker run -it -d \
  --name=order -p 8083:8080 \
  --net=solarmesh \
  -e APP_NAME=order \
  -e APP_VERSION=v1 \
  -e CLUSTER_NAME=cluster01 \
  -e UPSTREAM_URLS='http://product:8080/,http://inventory:8080/,http://payment:8080/,' \
  -e LOGGING_LEVEL=debug \
  solarmesh/solar-mock-app:latest
