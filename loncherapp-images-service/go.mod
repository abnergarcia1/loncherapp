module bitbucket.org/edgelabsolutions/loncherapp-images-service

replace bitbucket.org/edgelabsolutions/loncherapp-core => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-core

replace bitbucket.org/edgelabsolutions/loncherapp-protobuf => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-protobuf

go 1.14

require (
	bitbucket.org/edgelabsolutions/loncherapp-core v0.0.0-00010101000000-000000000000
	bitbucket.org/edgelabsolutions/loncherapp-protobuf v0.0.0-00010101000000-000000000000 // indirect
	github.com/aws/aws-sdk-go v1.37.1
	github.com/gin-gonic/gin v1.6.3
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/grpc v1.35.0 // indirect
)
