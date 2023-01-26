package grpc_client

import (
	"fmt"

	"gitlab.com/telegram_clone/api_gateway/config"
	pbc "gitlab.com/telegram_clone/api_gateway/genproto/chat_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:generate mockgen -source ../../genproto/chat_service/user_service_grpc.pb.go -package mock_grpc -destination ./mock_grpc/user_service_grpc.gen.go
//go:generate mockgen -source ../../genproto/chat_service/auth_service_grpc.pb.go -package mock_grpc -destination ./mock_grpc/auth_service_grpc.gen.go
//go:generate mockgen -source ../../genproto/chat_service/chat_service_grpc.pb.go -package mock_grpc -destination ./mock_grpc/chat_service_grpc.gen.go
//go:generate mockgen -source ../../genproto/chat_service/chat_message_service_grpc.pb.go -package mock_grpc -destination ./mock_grpc/message_service_grpc.gen.go

type GrpcClientI interface {
	UserService() pbc.UserServiceClient
	SetUserService(u pbc.UserServiceClient)
	AuthService() pbc.AuthServiceClient
	SetAuthService(u pbc.AuthServiceClient)
	ChatService() pbc.ChatServiceClient
	MessageService() pbc.MessageServiceClient
	SetChatService(p pbc.ChatServiceClient)
	SetMessageService(p pbc.MessageServiceClient)
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
		return nil, fmt.Errorf("chat service dial host: %s port:%s err: %v",
			cfg.ChatServiceHost, cfg.ChatServiceGrpcPort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user_service": pbc.NewUserServiceClient(connChatService),
			"auth_service": pbc.NewAuthServiceClient(connChatService),
			"chat_service": pbc.NewChatServiceClient(connChatService),
			"message_service": pbc.NewMessageServiceClient(connChatService),
		},
	}, nil
}

func (g *GrpcClient) UserService() pbc.UserServiceClient {
	return g.connections["user_service"].(pbc.UserServiceClient)
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