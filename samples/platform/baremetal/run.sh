#!/usr/bin/env bash

profile=dev

rm -f ./mockservices

if [[  $(ps aux | grep -c mockservices) -gt 1 ]]; then
  killall mockservices
fi

if [[ "$1" != "cleanup" ]]; then

  if [[ ! -f mockservices ]]; then
    echo "Assume that you have Go installed, run go build"
    go build
  fi

#    # Run as product v1
  nohup ./mockservices \
    --app.profiles.active=${profile} \
    --app.name=mysql \
    --app.version=v1 \
    --cluster.name=cluster01 \
    --user.data=baremetal \
    --tcp.server.enabled=true \
    --tcp.server.port=3306 \
    --grpc.server.port=7070 \
    --server.port=9090 &> mysql.out &

  # Run as product v1
  nohup ./mockservices \
    --app.profiles.active=${profile} \
    --app.name=product \
    --app.version=v1 \
    --cluster.name=cluster01 \
    --user.data=baremetal \
    --server.port=8080 \
    --grpc.server.port=7575 \
    --upstream.urls='tcp://localhost:3306,' &> product.out &


  # Run as inventory v1
  nohup ./mockservices \
    --app.profiles.active=${profile} \
    --app.name=inventory \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8081 \
    --grpc.server.port=7576 \
    --upstream.urls='grpc://localhost:7575,' &> inventory.out &


  # Run as order v2
  nohup ./mockservices \
    --app.profiles.active=${profile} \
    --app.name=order \
    --app.version=v2 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8082 \
    --grpc.server.port=7577 \
    --upstream.urls='grpc://localhost:7575,grpc://localhost:7576' &> order.out &


  # Run as payment v1
  nohup ./mockservices \
    --app.profiles.active=${profile} \
    --app.name=payment \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8083 \
    --grpc.server.port=7578 \
    --upstream.urls='http://localhost:8081/,http://localhost:8082/' &> payment.out &

fi
