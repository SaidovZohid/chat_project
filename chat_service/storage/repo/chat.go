package repo

import "time"

type ChatStrogeI interface {
	Create(chat *CreateChatReq) (*Chat, error)
	Get(id int64) (*Chat, error)
	Update(chat *Chat) (*Chat, error)
	Delete(id, user_id int64) error
	GetAll(params *GetAllChatsParams) (*GetAllChats, error)

	AddMember(*AddMemberRequest) error
	RemoveMember(*RemoveMemberRequest) error
	GetChatMembers(params *GetChatMembersParams) (*GetAllUsersResult, error)
}

type Chat struct {
	ID       int64
	Name     string
	UserID   int64
	UserInfo *GetUserInfo
	ChatType string
	ImageUrl string
}

type CreateChatReq struct {
	Name     string
	UserID   int64
	ChatType string
	ImageUrl string
	Members  []int64
}

type GetUserInfo struct {
	FirstName string
	LastName  string
	Email     string
	UserName  string
	ImageUrl  string
	CreatedAt time.Time
}

type GetAllChatsParams struct {
	Limit  int64
	Page   int64
	UserID int64
}

type GetAllChats struct {
	Chats []*Chat
	Count int64
}

type AddMemberRequest struct {
	ChatId int64
	UserId int64
}

type RemoveMemberRequest struct {
	ChatId int64
	UserId int64
}

type GetChatMembersParams struct {
	Limit  int64
	Page   int64
	ChatID int64
}
