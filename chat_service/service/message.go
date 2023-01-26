package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	pb "gitlab.com/telegram_clone/chat_service/genproto/chat_service"
	"gitlab.com/telegram_clone/chat_service/storage/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sirupsen/logrus"
	"gitlab.com/telegram_clone/chat_service/storage"
)

type MessageService struct {
	pb.UnimplementedMessageServiceServer
	storage storage.StorageI
	logger  *logrus.Logger
}

func NewMessageService(strg storage.StorageI, logger *logrus.Logger) *MessageService {
	return &MessageService{
		storage: strg,
		logger:  logger,
	}
}

func (s *MessageService) Create(ctx context.Context, req *pb.ChatMessage) (*pb.ChatMessage, error) {
	s.logger.Info("create message")
	chat, err := s.storage.ChatMessage().Create(&repo.ChatMessage{
		Message: req.Message,
		UserId:  req.UserId,
		ChatId:  req.ChatId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create message")
		return nil, status.Errorf(codes.Internal, "failed to create: %v", err)
	}

	return parseMessageModel(chat), nil
}

func (s *MessageService) Update(ctx context.Context, req *pb.ChatMessage) (*pb.ChatMessage, error) {
	chat, err := s.storage.ChatMessage().Update(&repo.ChatMessage{
		ID: req.Id,
		Message: req.Message,
		UserId: req.UserId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update message")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to update: %v", err)
	}

	return parseMessageModel(chat), nil
}

func (s *MessageService) Delete(ctx context.Context, req *pb.ChatIdRequest) (*emptypb.Empty, error) {
	err := s.storage.ChatMessage().Delete(req.Id, req.UserId)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete message")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *MessageService) GetAll(ctx context.Context, req *pb.GetAllMessagesParams) (*pb.GetAllMessages, error) {
	messages, err := s.storage.ChatMessage().GetAll(&repo.GetAllMessagesParams{
		Limit:  req.Limit,
		Page:   req.Page,
		ChatId: req.ChatId,
	})

	if err != nil {
		s.logger.WithError(err).Error("failed to get all messages")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get all: %w", err)
	}

	response := pb.GetAllMessages{
		Messages: make([]*pb.ChatMessage, 0),
		Count:    messages.Count,
	}
	for _, v := range messages.Messages {
		res := parseMessageModel(v)
		response.Messages = append(response.Messages, res)
	}

	return &response, nil
}

func parseMessageModel(res *repo.ChatMessage) *pb.ChatMessage {
	return &pb.ChatMessage{
		Id:      res.ID,
		Message: res.Message,
		UserId:  res.UserId,
		UserInfo: &pb.GetUserInfo{
			FirstName: res.UserInfo.FirstName,
			LastName:  res.UserInfo.LastName,
			Email:     res.UserInfo.Email,
			Username:  res.UserInfo.UserName,
			ImageUrl:  res.UserInfo.ImageUrl,
			CreatedAt: res.UserInfo.CreatedAt.Format(time.RFC3339),
		},
		ChatId:    res.ChatId,
		CreatedAt: res.CreatedAt.Format(time.RFC3339),
	}
}
