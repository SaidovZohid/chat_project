package main

import (
	"log"

	_ "github.com/lib/pq"

	"gitlab.com/telegram_clone/websocket_service/config"
	"gitlab.com/telegram_clone/websocket_service/websocket"

	grpcPkg "gitlab.com/telegram_clone/websocket_service/pkg/grpc_client"
)

func main() {
	cfg := config.Load(".")

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v", err)
	}

	websocket.Run(cfg, grpcConn)
}
