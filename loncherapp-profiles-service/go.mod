module bitbucket.org/edgelabsolutions/loncherapp-profiles-service

replace bitbucket.org/edgelabsolutions/loncherapp-core => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-core

replace bitbucket.org/edgelabsolutions/loncherapp-protobuf => /Users/abnergarcia/src/bitbucket.org/edgelab/loncherapp-protobuf

go 1.14

require (
	bitbucket.org/edgelabsolutions/loncherapp-core v0.2.1
	bitbucket.org/edgelabsolutions/loncherapp-protobuf v0.2.2
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/grpc v1.35.0
)
