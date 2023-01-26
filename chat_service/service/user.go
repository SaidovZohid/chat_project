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

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
	logger   *logrus.Logger
}

func NewUserService(strg storage.StorageI, inMemory storage.InMemoryStorageI, logger *logrus.Logger) *UserService {
	return &UserService{
		storage:  strg,
		inMemory: inMemory,
		logger:   logger,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	s.logger.Info("create user")
	user, err := s.storage.User().Create(&repo.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Password:        req.Password,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create user")
		return nil, status.Errorf(codes.Internal, "failed to create: %v", err)
	}

	return parseUserModel(user), nil
}

func (s *UserService) Get(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.storage.User().Get(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get: %v", err)
	}

	return parseUserModel(user), nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *pb.GetByEmailRequest) (*pb.User, error) {
	user, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user by email")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get by email: %v", err)
	}

	return parseUserModel(user), nil
}

func (s *UserService) GetAll(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	result, err := s.storage.User().GetAll(&repo.GetAllUsersParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get all user")
		return nil, status.Errorf(codes.Internal, "failed to get all users: %v", err)
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

func parseUserModel(user *repo.User) *pb.User {
	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}
}

func (s *UserService) Update(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Update(&repo.User{
		ID:              req.Id,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update user")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to update: %v", err)
	}

	return parseUserModel(user), nil
}

func (s *UserService) Delete(ctx context.Context, req *pb.GetUserRequest) (*emptypb.Empty, error) {
	err := s.storage.User().Delete(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete user")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) SetUserImage(ctx context.Context, req *pb.SetUserImageRequest) (*pb.User, error) {
	user, err := s.storage.User().SetUserImage(&repo.SetUserImageRequest{
		UserId:   req.UserId,
		ImageUrl: req.ImageUrl,
	})
	if err != nil {
		return &pb.User{}, err
	}

	return parseUserModel(user), nil
}
