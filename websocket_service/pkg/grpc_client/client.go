package grpc_client

import (
	"fmt"

	"gitlab.com/telegram_clone/websocket_service/config"
	pbc "gitlab.com/telegram_clone/websocket_service/genproto/chat_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	AuthService() pbc.AuthServiceClient
	ChatService() pbc.ChatServiceClient
	MessageService() pbc.MessageServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (GrpcClientI, error) {
	connChatService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ChatServiceHost, cfg.ChatServiceGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %v",
			cfg.ChatServiceHost, cfg.ChatServiceGrpcPort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"auth_service":    pbc.NewAuthServiceClient(connChatService),
			"chat_service":    pbc.NewChatServiceClient(connChatService),
			"message_service": pbc.NewMessageServiceClient(connChatService),
		},
	}, nil
}

func (g *GrpcClient) AuthService() pbc.AuthServiceClient {
	return g.connections["auth_service"].(pbc.AuthServiceClient)
}

func (g *GrpcClient) ChatService() pbc.ChatServiceClient {
	return g.connections["chat_service"].(pbc.ChatServiceClient)
}

func (g *GrpcClient) MessageService() pbc.MessageServiceClient {
	return g.connections["message_service"].(pbc.MessageServiceClient)
}
