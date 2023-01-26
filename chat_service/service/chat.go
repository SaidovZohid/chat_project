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

type ChatService struct {
	pb.UnimplementedChatServiceServer
	storage storage.StorageI
	logger  *logrus.Logger
}

func NewChatService(strg storage.StorageI, logger *logrus.Logger) *ChatService {
	return &ChatService{
		storage: strg,
		logger:  logger,
	}
}

func (s *ChatService) Create(ctx context.Context, req *pb.CreateChatReq) (*pb.Chat, error) {
	s.logger.Info("create chat")
	chat, err := s.storage.Chat().Create(&repo.CreateChatReq{
		Name:     req.Name,
		UserID:   req.UserId,
		ChatType: req.ChatType,
		ImageUrl: req.ImageUrl,
		Members:  req.Members,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create chat")
		return nil, status.Errorf(codes.Internal, "failed to create: %v", err)
	}

	return parseChatModel(chat), nil
}

func parseChatModel(res *repo.Chat) *pb.Chat {
	return &pb.Chat{
		Id:     res.ID,
		Name:   res.Name,
		UserId: res.UserID,
		UserInfo: &pb.GetUserInfo{
			FirstName: res.UserInfo.FirstName,
			LastName:  res.UserInfo.LastName,
			Email:     res.UserInfo.Email,
			Username:  res.UserInfo.UserName,
			ImageUrl:  res.UserInfo.ImageUrl,
			CreatedAt: res.UserInfo.CreatedAt.Format(time.RFC3339),
		},
		ChatType: res.ChatType,
		ImageUrl: res.ImageUrl,
	}
}

func (s *ChatService) Get(ctx context.Context, req *pb.IdRequest) (*pb.Chat, error) {
	chat, err := s.storage.Chat().Get(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to get chat info")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get: %v", err)
	}

	return parseChatModel(chat), nil
}

func (s *ChatService) Update(ctx context.Context, req *pb.Chat) (*pb.Chat, error) {
	chat, err := s.storage.Chat().Update(&repo.Chat{
		ID:       req.Id,
		Name:     req.Name,
		UserID:   req.UserId,
		ImageUrl: req.ImageUrl,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update chat")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get: %v", err)
	}

	return parseChatModel(chat), nil
}

func (s *ChatService) Delete(ctx context.Context, req *pb.ChatIdRequest) (*emptypb.Empty, error) {
	err := s.storage.Chat().Delete(req.Id, req.UserId)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete chat")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ChatService) GetAll(ctx context.Context, req *pb.GetAllChatsParams) (*pb.GetAllChatsRes, error) {
	privateChats, err := s.storage.Chat().GetAll(&repo.GetAllChatsParams{
		Limit:  req.Limit,
		Page:   req.Page,
		UserID: req.UserId,
	})

	if err != nil {
		s.logger.WithError(err).Error("failed to get all private chats")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get all: %w", err)
	}

	response := pb.GetAllChatsRes{
		Chats: make([]*pb.Chat, 0),
		Count: privateChats.Count,
	}
	for _, v := range privateChats.Chats {
		res := parseChatModel(v)
		response.Chats = append(response.Chats, res)
	}

	return &response, nil
}

// Group chat methods
func (s *ChatService) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*emptypb.Empty, error) {
	s.logger.Info("Add members")
	err := s.storage.Chat().AddMember(&repo.AddMemberRequest{
		ChatId: req.ChatId,
		UserId: req.UserId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to add member")
		return nil, status.Errorf(codes.Internal, "failed to addmember: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ChatService) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*emptypb.Empty, error) {
	s.logger.Info("Remove members")
	err := s.storage.Chat().RemoveMember(&repo.RemoveMemberRequest{
		ChatId: req.ChatId,
		UserId: req.UserId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to remove member")
		return nil, status.Errorf(codes.Internal, "failed to remove member: %v", err)
	}
	return nil, nil
}

func (s *ChatService) GetChatMembers(ctx context.Context, req *pb.GetChatMembersParams) (*pb.GetAllUsersResponse, error) {
	result, err := s.storage.Chat().GetChatMembers(&repo.GetChatMembersParams{
		Limit:  req.Limit,
		Page:   req.Page,
		ChatID: req.ChatId,
	})

	if err != nil {
		s.logger.WithError(err).Error("failed to get all private chats")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get all: %w", err)
	}

	response := pb.GetAllUsersResponse{
		Count: result.Count,
		Users: make([]*pb.User, 0),
	}

	for _, user := range result.Users {
		response.Users = append(response.Users, parseUserModel(user))
	}

	return &response, nil
}
