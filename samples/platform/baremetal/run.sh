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
    --upstream.urls='http://localhost:8080/,' &


  # Run as order v2
  nohup ./mockservices \
    --app.name=order \
    --app.version=v2 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8082 \
    --upstream.urls='http://localhost:8080/,http://localhost:8081/' &


  # Run as payment v1
  nohup ./mockservices \
    --app.name=payment \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8083 \
    --upstream.urls='http://localhost:8081/,http://localhost:8082/' &

fi
