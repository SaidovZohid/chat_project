package repo

import "time"

type ChatMessageStrogeI interface {
	Create(m *ChatMessage) (*ChatMessage, error)
	Update(m *ChatMessage) (*ChatMessage, error)
	Delete(id, user_id int64) error
	GetAll(params *GetAllMessagesParams) (*GetAllMessages, error)
}

type ChatMessage struct {
	ID        int64
	Message   string
	UserId    int64
	UserInfo  *GetUserInfo
	ChatId    int64
	CreatedAt time.Time
}

type GetAllMessagesParams struct {
	Limit  int64
	Page   int64
	ChatId int64
}

type GetAllMessages struct {
	Messages []*ChatMessage
	Count    int64
}
