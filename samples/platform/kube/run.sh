#!/usr/bin/env bash

#echo "use shell command to produce yaml that can be deployed to kubernetes"
function new_app() {
IMAGE="$1"
APP_PROFILES_ACTIVE="$2"
APP_NAME="$3"
APP_VERSION="$4"
CLUSTER_NAME="$5"
UPSTREAM_URLS="$6"
NODE_PORT=$7

if [[ "${NODE_PORT}" != "" ]]; then

kubectl apply -f - <<_EOF
apiVersion: v1
kind: Service
metadata:
  name: ${APP_NAME}
  labels:
    app: ${APP_NAME}
spec:
  ports:
    - port: 8080
      name: http
      targetPort: 8080
      nodePort: ${NODE_PORT}
    - port: 7575
      name: grpc
      targetPort: 7575
    - port: 8585
      name: tcp
      targetPort: 8585
  selector:
    app: ${APP_NAME}
  type: NodePort
_EOF

else

kubectl apply -f - <<_EOF
apiVersion: v1
kind: Service
metadata:
  name: ${APP_NAME}
  labels:
    app: ${APP_NAME}
spec:
  ports:
    - port: 8080
      name: http
      targetPort: 8080
    - port: 7575
      name: grpc
      targetPort: 7575
    - port: 8585
      name: tcp
      targetPort: 8585
  selector:
    app: ${APP_NAME}
  type: ClusterIP
_EOF
fi


kubectl apply -f - <<_EOF

apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${APP_NAME}-${APP_VERSION}
  labels:
    app: ${APP_NAME}
    version: ${APP_VERSION}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${APP_NAME}
      version: ${APP_VERSION}
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "true"
      labels:
        app: ${APP_NAME}
        version: ${APP_VERSION}
    spec:
      terminationGracePeriodSeconds: 0
      containers:
        - name: ${APP_NAME}
          image: ${IMAGE}
          imagePullPolicy: Always # IfNotPresent
          env:
            - name: APP_PROFILES_ACTIVE
              value: "${APP_PROFILES_ACTIVE:${profile}}"
            - name: APP_NAME
              value: "${APP_NAME}"
            - name: APP_VERSION
              value: "${APP_VERSION}"
            - name: CLUSTER_NAME
              value: "${CLUSTER_NAME}"
            - name: USER_DATA
              value: "${USER_DATA:Your own data}"
            - name: TCP_SERVER_ENABLED
              value: "true"
            - name: UPSTREAM_URLS
              value: "${UPSTREAM_URLS}"

          ports:
            - containerPort: 8080
            - containerPort: 7575
            - containerPort: 8585
_EOF
}

node_port=32080
if [[ "$3" != "" ]]; then
  node_port=$3
fi

image=solarmesh/mockservices:latest
if [[ "$2" != "" ]]; then
  image=$2
fi

profile=dev
if [[ "$1" != "" ]]; then
  profile=$1
fi

gateway_upstreams="http://payment:8080,http://order:8080,http://user:8080,http://reviews:8080,http://recommendation:8080,http://category:8080,"

new_app ${image} ${profile} api-gateway v1 cluster01 "${gateway_upstreams}" "${node_port}"

new_app ${image} ${profile} order v1 cluster01 "http://payment:8080,http://user:8080,http://reviews:8080,http://inventory:8080,"
new_app ${image} ${profile} order v2 cluster01 "http://payment:8080,http://user:8080,http://reviews:8080,http://inventory:8080,"

new_app ${image} ${profile} payment v1 cluster01 "http://user:8080,http://payment-db:8080"
new_app ${image} ${profile} payment v2 cluster01 "http://user:8080,http://payment-db:8080"

new_app ${image} ${profile} payment-db v1 cluster01 ""

new_app ${image} ${profile} inventory v1 cluster01 "http://user:8080,http://product:8080,"

new_app ${image} ${profile} reviews v1 cluster01 "http://user:8080,"
new_app ${image} ${profile} reviews v2 cluster01 "http://user:8080,"
new_app ${image} ${profile} reviews v3 cluster01 "http://user:8080,"

new_app ${image} ${profile} recommendation v1 cluster01 "http://user:8080,"

new_app ${image} ${profile} category v1 cluster01 "http://product:8080,"
new_app ${image} ${profile} category v2 cluster01 "grpc://product:7575,"

new_app ${image} ${profile} product v1 cluster01 "http://product-db:8080,"
new_app ${image} ${profile} product v2 cluster01 "http://product-db:8080,"

new_app ${image} ${profile} product-db v1 cluster01 ""

new_app ${image} ${profile} user v1 cluster01 "http://user-db:8080,"
new_app ${image} ${profile} user-db v1 cluster01 ""
