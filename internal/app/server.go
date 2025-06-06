package app

import (
	"fmt"
	"github.com/GalahadKingsman/messenger_dialog/internal/app/dialogservice"
	"github.com/GalahadKingsman/messenger_dialog/internal/config"
	"github.com/GalahadKingsman/messenger_dialog/internal/database"
	"github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(config *config.Config) {
	db, err := database.Init(config.DB)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	userRepo := user_repo.New(db)
	service := dialogservice.New(userRepo)
	messenger_dialog_api.RegisterUserServiceServer(grpcServer, service)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Println(err)
	}
}
