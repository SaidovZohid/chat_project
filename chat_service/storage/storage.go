package storage

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/telegram_clone/chat_service/storage/postgres"
	"gitlab.com/telegram_clone/chat_service/storage/repo"
)

type StorageI interface {
	User() repo.UserStorageI
	Permission() repo.PermissionStorageI
	Chat() repo.ChatStrogeI
	ChatMessage() repo.ChatMessageStrogeI
}

type storagePg struct {
	userRepo        repo.UserStorageI
	permissionRepo  repo.PermissionStorageI
	chatRepo        repo.ChatStrogeI
	chatMessageRepo repo.ChatMessageStrogeI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		userRepo:        postgres.NewUser(db),
		permissionRepo:  postgres.NewPermission(db),
		chatRepo:        postgres.NewChat(db),
		chatMessageRepo: postgres.NewChatMessage(db),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *storagePg) Permission() repo.PermissionStorageI {
	return s.permissionRepo
}

func (s *storagePg) Chat() repo.ChatStrogeI {
	return s.chatRepo
}

func (s *storagePg) ChatMessage() repo.ChatMessageStrogeI {
	return s.chatMessageRepo
}
