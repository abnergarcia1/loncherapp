    export ACCESS_SECRET=""
    export REDIS_DSN="localhost:6380"
    export REDIS_PASSWORD=""
    export REFRESH_SECRET="mcmvmkmsdnfsdmfdsjf"
    export IMAGE_SERVICE_HOST="localhost:8082"
    export LONCHERAPP_DB_CONNECTION="root:password1234@tcp(localhost:3307)/loncherapp?parseTime=true"

    #aws
    export AWS_REGION="us-west-1"
    export BUCKET_NAME="loncherapp-content"
    export AWS_ACCESS_KEY_ID=""
    export AWS_SECRET_ACCESS_KEY=""


  	go run ./app/main.go