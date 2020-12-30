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


# Run as order v2
./solar-mock-app --app.name=order --app.version=v2 --cluster.name=cluster02 --user.data=demo --server.port=8082 --upstream.urls=http://localhost:8080/,http://localhost:8081/


# Run as payment v1
./solar-mock-app --app.name=payment --app.version=v1 --cluster.name=cluster02 --user.data=demo --server.port=8083 --upstream.urls=http://localhost:8080/,http://localhost:8081/,http://localhost:8082/

```

### Run with Docker image


```bash

docker run -it --name 

```