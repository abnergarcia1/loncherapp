#!/bin/bash

    export LONCHERAPP_DB_CONNECTION="root:password1234@tcp(localhost:3307)/loncherapp?parseTime=true"
    #"loncherapp-admin:L0nCh3r@pP@tcp(mariadb-18224-0.cloudclusters.net:18224)/loncherapp?parseTime=true"




  	go run ./app/main.go