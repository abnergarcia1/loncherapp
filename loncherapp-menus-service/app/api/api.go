package api

import (
	"context"
	"net"
	"os"

	pb "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/menus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"

	"google.golang.org/grpc/reflection"

	handlers "bitbucket.org/edgelabsolutions/loncherapp-menus-service/app/handlers"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	defers []func()
	port   string = ":7550"
)

func StartAPI() (close func(), err error) {

	ctx := context.Background()

	sql.Init(
		sql.SetConnectionString(os.Getenv("LONCHERAPP_DB_CONNECTION")),
	)

	server := handlers.NewMenusAPIServer(ctx)

	grpcServer := grpc.NewServer()

	pb.RegisterMenusServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Info("GRPC Menus started and serving")

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
