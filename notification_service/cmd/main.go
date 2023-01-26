package main

import (
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gitlab.com/telegram_clone/notification_service/genproto/notification_service"

	"gitlab.com/telegram_clone/notification_service/config"
	"gitlab.com/telegram_clone/notification_service/service"
)

func main() {
	cfg := config.Load(".")

	notificationService := service.NewNotificationService(&cfg)

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterNotificationServiceServer(s, notificationService)

	log.Println("Grpc server started in port ", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}

}
