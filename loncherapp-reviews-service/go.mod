module bitbucket.org/edgelabsolutions/loncherapp-reviews-service

replace bitbucket.org/edgelabsolutions/loncherapp-core => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-core

replace bitbucket.org/edgelabsolutions/loncherapp-protobuf => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-protobuf

go 1.14

require (
	bitbucket.org/edgelabsolutions/loncherapp-core v0.0.0-00010101000000-000000000000
	bitbucket.org/edgelabsolutions/loncherapp-protobuf v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.0
	google.golang.org/grpc v1.36.0
)
