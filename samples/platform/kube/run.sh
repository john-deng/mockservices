#!/usr/bin/env bash

BASEDIR=$(dirname "$0")

# shellcheck disable=SC1090
source "${BASEDIR}/new"

export APP_PROFILES_ACTIVE=dev
if [[ "$1" != "" ]]; then
  export APP_PROFILES_ACTIVE=$1
fi

export IMAGE=solarmesh/mockservices:latest
if [[ "$2" != "" ]]; then
  export IMAGE=$2
fi

GW_NODE_PORT=32080
if [[ "$3" != "" ]]; then
  GW_NODE_PORT=$3
fi

gateway_upstream_urls="http://payment:8080,http://order:8080,http://user:8080,http://reviews:8080,http://recommendation:8080,http://category:8080,"

new_app api-gateway v1 cluster01 "${gateway_upstream_urls}" "${GW_NODE_PORT}"

order_upstream_urls="grpc://cart:7575,grpc://payment:7575,grpc://user:7575,grpc://reviews:7575,grpc://inventory:7575,grpc://logistics:7575,grpc://shipment:7575,grpc://notification:7575"

new_app order v1 cluster01 "${order_upstream_urls}"
new_app order v2 cluster01 "${order_upstream_urls}"
new_app order v3 cluster01 "${order_upstream_urls}"

new_app cart v1 cluster01 "tcp://cart-db:8585"
new_app cart-db v1 cluster01 ""

new_app payment v1 cluster01 "grpc://user:7575,grpc://notification:7575,tcp://payment-db:8585"

new_app logistics v1 cluster01 "grpc://user:7575,grpc://notification:7575,"

new_app shipment v1 cluster01 "grpc://user:7575,grpc://notification:7575,"

new_app payment-db v1 cluster01 ""

new_app inventory v1 cluster01 "grpc://user:7575,grpc://product:7575,grpc://notification:7575,"

new_app reviews v1 cluster01 "grpc://user:7575,grpc://notification:7575,"
new_app reviews v2 cluster01 "grpc://user:7575,grpc://notification:7575,"
new_app reviews v3 cluster01 "grpc://user:7575,grpc://notification:7575,"

new_app recommendation v1 cluster01 "grpc://user:7575,grpc://notification:7575,"

new_app category v1 cluster01 "http://product:8080,grpc://notification:7575,"
new_app category v2 cluster01 "grpc://product:7575,grpc://notification:7575,"

new_app product v1 cluster01 "tcp://product-db:8585,grpc://notification:7575,"

new_app product-db v1 cluster01 ""

new_app notification v1 cluster01 "tcp://kafka:8585,"
new_app kafka v1 cluster01 ""

new_app logs v1 cluster01 ""

new_app user v1 cluster01 "grpc://logs:7575,tcp://user-db:8585,"
new_app user-db v1 cluster01 ""
