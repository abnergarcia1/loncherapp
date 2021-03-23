package api

import (
	"context"
	"net"
	"os"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	"google.golang.org/grpc/reflection"

	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/users"
	handlers "bitbucket.org/edgelabsolutions/loncherapp-users-service/app/grpc-handlers"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	defers []func()
	port   string = ":7510"
)

func StartAPI() (close func(), err error) {

	//os.Setenv("LONCHERAPP_DB_CONNECTION", "loncherapp-admin:L0nCh3r@pP@tcp(mariadb-18224-0.cloudclusters.net:18224)/loncherapp?parseTime=true")

	ctx := context.Background()

	sql.Init(
		sql.SetConnectionString(os.Getenv("LONCHERAPP_DB_CONNECTION")),
	)

	server := handlers.NewUsersAPIServer(ctx)

	grpcServer := grpc.NewServer()

	pb.RegisterUsersServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Info("GRPC Users started and serving")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen :%v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}

	return func() {
		for _, c := range defers {
			c()
			grpcServer.GracefulStop()
		}
	}, nil

}
