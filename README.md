# solar-mock-app

solar-mock-app is an all in one app that build for testing Service Mesh, the key feature of this app is that you can deploy as many apps as possible with just a single binary file or docker image.

## Getting started

### From Source Code

Get the source code and build with Go tools

```bash
# Clone the source code
git clone https://github.com/solarmesh-io/solar-mock-app.git

# Build
cd solar-mock-app
go build

```

After the app is built, you can run it directly

### 1. Run the binary file.
```bash

# Run as product v1
./solar-mock-app --app.name=product --app.version=v1 --cluster.name=cluster01 --user.data=demo --server.port=8080 


# Run as inventory v1
./solar-mock-app --app.name=inventory --app.version=v1 --cluster.name=cluster02 --user.data=demo --server.port=8081 --upstream.urls=http://localhost:8080/,


# Run as payment v2
./solar-mock-app --app.name=payment --app.version=v2 --cluster.name=cluster02 --user.data=demo --server.port=8082 --upstream.urls=http://localhost:8080/,http://localhost:8081/


# Run as order v1
./solar-mock-app --app.name=order --app.version=v1 --cluster.name=cluster02 --user.data=demo --server.port=8083 --upstream.urls=http://localhost:8080/,http://localhost:8081/,http://localhost:8082/

```

### Run with Docker image


```bash

docker network create solarmesh 

# Run as product v1
docker run -it -d \
  --name=product \
  --net=solarmesh \
  -p 8080:8080 \
  -e APP_NAME=product \
  -e APP_VERSION=v1 \
  -e CLUSTER_NAME=cluster01 \
  solarmesh/solar-mock-app:latest


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

# Run as payment v2
docker run -it -d \
  --name=payment \
  --net=solarmesh \
  -p 8082:8080 \
  -e APP_NAME=payment \
  -e APP_VERSION=v2 \
  -e CLUSTER_NAME=cluster01 \
  -e UPSTREAM_URLS='http://product:8080/,http://product:8080/,' \
  solarmesh/solar-mock-app:latest


# Run as order v1
docker run -it -d \
  --name=order -p 8083:8080 \
  --net=solarmesh \
  -e APP_NAME=order \
  -e APP_VERSION=v1 \
  -e CLUSTER_NAME=cluster01 \
  -e UPSTREAM_URLS='http://product:8080/,http://product:8080/,http://payment:8080/,' \
  -e LOGGING_LEVEL=debug \
  solarmesh/solar-mock-app:latest

```