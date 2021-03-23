#!/bin/bash

    export ACCESS_SECRET="keysecret"
    export GATEWAY_PORT="8081"
    export REDIS_DSN="localhost:6380"
    export REDIS_PASSWORD=""
    export REFRESH_SECRET="mcmvmkmsdnfsdmfdsjf"
    export USER_SERVICE_HOST="localhost:7510"
    export PROFILES_SERVICE_HOST="localhost:7530"
    export ROUTES_SERVICE_HOST="localhost:7540"
    export MENUS_SERVICE_HOST="localhost:7550"
    export REVIEWS_SERVICE_HOST="localhost:7560"


  #  sleep 10
  	go run ./app/main.go