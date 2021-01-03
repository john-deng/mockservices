#!/usr/bin/env bash

rm -f ./mockservices


if [[  $(ps aux | grep -c mockservices) -gt 1 ]]; then
  killall mockservices
fi

if [[ "$1" != "cleanup" ]]; then

  if [[ ! -f mockservices ]]; then
    echo "Assume that you have Go installed, run go build"
    go build
  fi

  # Run as product v1
  nohup ./mockservices \
    --app.name=product \
    --app.version=v1 \
    --cluster.name=cluster01 \
    --user.data=baremetal \
    --server.port=8080 &


  # Run as inventory v1
  nohup ./mockservices \
    --app.name=inventory \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8081 \
    --grpc.server.port=7576 \
    --upstream.urls='grpc://localhost:7575,' &


  # Run as order v2
  nohup ./mockservices \
    --app.name=order \
    --app.version=v2 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8082 \
    --grpc.server.port=7577 \
    --upstream.urls='grpc://localhost:7575,grpc://localhost:7576' &


  # Run as payment v1
  nohup ./mockservices \
    --app.name=payment \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8083 \
    --grpc.server.port=7578 \
    --upstream.urls='http://localhost:8081/,http://localhost:8082/' &

fi
