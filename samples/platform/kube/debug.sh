#!/usr/bin/env bash

BASEDIR=$(dirname "$0")

# shellcheck disable=SC1090
source "${BASEDIR}/new"

export APP_PROFILES_ACTIVE=debug

export IMAGE=solarmesh/ide-go:v2.6.3
if [[ "$1" != "" ]]; then
  export IMAGE=$1
fi

new_app user v1 cluster01 "http://user-db:8080,grpc://user-db:7575,tcp://user-db:8585" "32180" "32143"
new_app user-db v1 cluster01 "" "32181" "32144"

