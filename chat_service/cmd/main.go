package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pb "gitlab.com/telegram_clone/chat_service/genproto/chat_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gitlab.com/telegram_clone/chat_service/config"
	"gitlab.com/telegram_clone/chat_service/service"
	"gitlab.com/telegram_clone/chat_service/storage"

	"gitlab.com/telegram_clone/chat_service/pkg/cronjob"
	grpcPkg "gitlab.com/telegram_clone/chat_service/pkg/grpc_client"
	"gitlab.com/telegram_clone/chat_service/pkg/logger"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v", err)
	}
	logrus := logger.New()

	// Registering jobs
	cron := cronjob.NewCronjob(strg, grpcConn, &cfg, logrus)
	cron.RegisterTasks()

	userService := service.NewUserService(strg, inMemory, logrus)
	authService := service.NewAuthService(strg, inMemory, grpcConn, &cfg, logrus)
	chatService := service.NewChatService(strg, logrus)
	messageService := service.NewMessageService(strg, logrus)

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterUserServiceServer(s, userService)
	pb.RegisterAuthServiceServer(s, authService)
	pb.RegisterChatServiceServer(s, chatService)
	pb.RegisterMessageServiceServer(s, messageService)

	log.Println("Grpc server started in port ", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}

}
