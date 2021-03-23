     export ACCESS_SECRET=""
    export REDIS_DSN="localhost:6380"
    export REDIS_PASSWORD=""
    export REFRESH_SECRET="mcmvmkmsdnfsdmfdsjf"
    export LONCHERAPP_DB_CONNECTION="root:password1234@tcp(localhost:3307)/loncherapp?parseTime=true"
    export PAYMENTS_SERVICE_HOST="localhost:8080"
    export PAYPAL_CLIENT_ID=""
    export PAYPAL_CLIENT_SECRET=""
    export PAYPAL_OAUTH_API="https://api-m.sandbox.paypal.com/v1/oauth2/token"
    export PAYPAL_ORDER_API="https://api-m.sandbox.paypal.com/v2/checkout/orders"
    export LONCHERAPP_MONGO_DB_URI="mongodb+srv:///myFirstDatabase?retryWrites=true&w=majority"

  	go run ./app/main.go