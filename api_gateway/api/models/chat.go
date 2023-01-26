package models

type Chat struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	UserID   int64       `json:"user_id"`
	UserInfo GetUserInfo `json:"user_info"`
	ChatType string      `json:"chat_type"`
	ImageUrl string      `json:"image_url"`
}

type ChatReq struct {
	Name     string  `json:"name" binding:"required"`
	ChatType string  `json:"chat_type" binding:"required"`
	ImageUrl string  `json:"image_url" binding:"required"`
	Members  []int64 `json:"members" binding:"required"`
}

type ChatRes struct {
	ID       int64       `json:"id"`
	UserID   int64       `json:"user_id"`
	UserInfo GetUserInfo `json:"user_info"`
}

type GetUserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	ImageUrl  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type GetAllChatsParams struct {
	Limit  int64 `json:"limit" binding:"required" default:"10"`
	Page   int64 `json:"page" binding:"required" default:"1"`
	UserID int64 `json:"user_id"`
}

type GetAllChatsRes struct {
	Chats []*Chat `json:"private_chats"`
	Count int64   `json:"count"`
}

type AddRemoveMemberReq struct {
	UserID int64 `json:"user_id" binding:"required"`
	ChatID int64 `json:"chat_id" binding:"required"`
}

type LeaveGroupReq struct {
	ChatID int64 `json:"chat_id" binding:"required"`
}

type GetChatMembersParams struct {
	Limit  int64 `json:"limit" binding:"required" default:"10"`
	Page   int64 `json:"page" binding:"required" default:"1"`
	ChatID int64 `json:"chat_id" binding:"required"`
}
