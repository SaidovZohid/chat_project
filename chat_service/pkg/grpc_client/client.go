package grpc_client

import (
	"fmt"

	"gitlab.com/telegram_clone/chat_service/config"
	pbn "gitlab.com/telegram_clone/chat_service/genproto/notification_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	NotificationService() pbn.NotificationServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (GrpcClientI, error) {
	connNotificationService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.NotificationServiceHost, cfg.NotificationServiceGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %v",
			cfg.NotificationServiceHost, cfg.NotificationServiceGrpcPort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"notification_service": pbn.NewNotificationServiceClient(connNotificationService),
		},
	}, nil
}

func (g *GrpcClient) NotificationService() pbn.NotificationServiceClient {
	return g.connections["notification_service"].(pbn.NotificationServiceClient)
}
