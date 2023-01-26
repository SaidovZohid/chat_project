package models

type Message struct {
	ID int64 `json:"id"`
	Message string `json:"message"`
	UserID int64 `json:"user_id"`
	UserInfo GetUserInfo `json:"user_info"`
	ChatID int64 `json:"chat_id"`
	CreatedAt string `json:"created_at"`
}

type MessageReq struct {
	Message string `json:"message" binding:"required"`
	UserID int64 `json:"user_id" binding:"required"`
	ChatID int64 `json:"chat_id" binding:"required"`
}

type GetAllMessagesRes struct {
	Messages []*Message `json:"messages" binding:"required"`
	Count int64 `json:"count" binding:"required"`
}

type GetAllMessagesParams struct {
	Limit int64 `json:"limit" binding:"required" default:"10"`
	Page  int64 `json:"page" binding:"required" default:"1"`
	ChatID int64  `json:"chat_id"`
}

