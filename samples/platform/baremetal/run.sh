#!/usr/bin/env bash

rm -f ./solar-mock-app

if [[ ! -f solar-mock-app ]]; then
  echo "Assume that you have Go installed, run go build"
  go build
fi

if [[  $(ps aux | grep -c solar-mock-app) -gt 1 ]]; then
  killall solar-mock-app
fi

if [[ "$1" != "cleanup" ]]; then
  # Run as product v1
  nohup ./solar-mock-app \
    --app.name=product \
    --app.version=v1 \
    --cluster.name=cluster01 \
    --user.data=baremetal \
    --server.port=8080 &


  # Run as inventory v1
  nohup ./solar-mock-app \
    --app.name=inventory \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8081 \
    --upstream.urls='http://localhost:8080/,' &


  # Run as order v2
  nohup ./solar-mock-app \
    --app.name=order \
    --app.version=v2 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8082 \
    --upstream.urls='http://localhost:8080/,http://localhost:8081/' &


  # Run as payment v1
  nohup ./solar-mock-app \
    --app.name=payment \
    --app.version=v1 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8083 \
    --upstream.urls='http://localhost:8080/,http://localhost:8081/,http://localhost:8082/' &

fi

./solar-mock-app \
    --app.name=payment \
    --app.version=v2 \
    --cluster.name=cluster02 \
    --user.data=demo \
    --server.port=8084 \
    --logging.level=debug \
    --upstream.urls='http://localhost:8080/,http://localhost:8081/,http://localhost:8082/'