package service

import (
	"context"

	"gitlab.com/telegram_clone/notification_service/config"
	pb "gitlab.com/telegram_clone/notification_service/genproto/notification_service"
	emailPkg "gitlab.com/telegram_clone/notification_service/pkg/email"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		cfg: cfg,
	}
}

func (s *NotificationService) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*emptypb.Empty, error) {
	err := emailPkg.SendEmail(s.cfg, &emailPkg.SendEmailRequest{
		To:      []string{req.To},
		Subject: req.Subject,
		Body:    req.Body,
		Type:    req.Type,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}
	return &emptypb.Empty{}, nil

}
