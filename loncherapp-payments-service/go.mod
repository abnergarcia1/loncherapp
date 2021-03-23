module bitbucket.org/edgelabsolutions/loncherapp-payments-service

replace bitbucket.org/edgelabsolutions/loncherapp-core => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-core

replace bitbucket.org/edgelabsolutions/loncherapp-protobuf => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-protobuf

go 1.14

require (
	bitbucket.org/edgelabsolutions/loncherapp-core v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.37.29 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	go.mongodb.org/mongo-driver v1.5.0
)
