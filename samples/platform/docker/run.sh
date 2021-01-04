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

clean_up_docker database
# Run as product v1
if [[ "$1" != "cleanup" ]]; then
  docker run -it -d \
    --name=database \
    --net=solarmesh \
    -e APP_NAME=database \
    -e APP_VERSION=v1 \
    -e CLUSTER_NAME=cluster01 \
    -e APP_PROFILES_ACTIVE=local \
    -e TCP_SERVER_ENABLED=true \
    solarmesh/mockservices:latest
fi

clean_up_docker product
# Run as product v1
if [[ "$1" != "cleanup" ]]; then
  docker run -it -d \
    --name=product \
    --net=solarmesh \
    -e APP_NAME=product \
    -e APP_VERSION=v1 \
    -e CLUSTER_NAME=cluster01 \
    -e APP_PROFILES_ACTIVE=local \
    -e UPSTREAM_URLS='tcp://database:8585,' \
    solarmesh/mockservices:latest
fi
clean_up_docker inventory
# Run as inventory v1
if [[ "$1" != "cleanup" ]]; then
  docker run -it -d \
    --name=inventory \
    --net=solarmesh \
    -e APP_NAME=inventory \
    -e APP_VERSION=v1 \
    -e CLUSTER_NAME=cluster01 \
    -e APP_PROFILES_ACTIVE=local \
    -e UPSTREAM_URLS='grpc://product:7575,' \
    solarmesh/mockservices:latest
fi
clean_up_docker payment
# Run as payment v2
if [[ "$1" != "cleanup" ]]; then
  docker run -it -d \
    --name=payment \
    --net=solarmesh \
    -e APP_NAME=payment \
    -e APP_VERSION=v2 \
    -e CLUSTER_NAME=cluster01 \
    -e APP_PROFILES_ACTIVE=local \
    -e UPSTREAM_URLS='grpc://product:7575,grpc://inventory:7575,' \
    solarmesh/mockservices:latest
fi
clean_up_docker order
# Run as order v1
if [[ "$1" != "cleanup" ]]; then
  docker run -it -d \
    --name=order \
    -p 8080:8080 \
    --net=solarmesh \
    -e APP_NAME=order \
    -e APP_VERSION=v1 \
    -e CLUSTER_NAME=cluster01 \
    -e APP_PROFILES_ACTIVE=local \
    -e UPSTREAM_URLS='http://inventory:8080/,http://payment:8080/,' \
    -e LOGGING_LEVEL=debug \
    solarmesh/mockservices:latest
fi
