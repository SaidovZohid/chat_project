package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/telegram_clone/chat_service/storage/repo"
)

type chatMessageRepo struct {
	db *sqlx.DB
}

func NewChatMessage(db *sqlx.DB) repo.ChatMessageStrogeI {
	return &chatMessageRepo{
		db: db,
	}
}

func (pr *chatMessageRepo) Create(message *repo.ChatMessage) (*repo.ChatMessage, error) {
	query := `
		INSERT INTO chat_messages (
			message,
			user_id,
			chat_id
		) VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	err := pr.db.QueryRow(
		query,
		message.Message,
		message.UserId,
		message.ChatId,
	).Scan(
		&message.ID,
		&message.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	message.UserInfo, err = getUserInfo(pr.db, message.UserId)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (pr *chatMessageRepo) Update(message *repo.ChatMessage) (*repo.ChatMessage, error) {
	query := `
		UPDATE chat_messages SET
			message = $1
		WHERE id=$2 and user_id=$3
		RETURNING
			chat_id,
			created_at
	`

	err := pr.db.QueryRow(
		query,
		message.Message,
		message.ID,
		message.UserId,
	).Scan(
		&message.ChatId,
		&message.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	message.UserInfo, err = getUserInfo(pr.db, message.UserId)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (pr *chatMessageRepo) Delete(id, user_id int64) error {
	query := " DELETE FROM chat_messages WHERE id=$1 and user_id = $2 "

	result, err := pr.db.Exec(query, id, user_id)
	if err != nil {
		return err
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pr *chatMessageRepo) GetAll(params *repo.GetAllMessagesParams) (*repo.GetAllMessages, error) {
	result := repo.GetAllMessages{
		Messages: make([]*repo.ChatMessage, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	filter := ""
	if params.ChatId > 0 {
		filter = fmt.Sprintf(" WHERE chat_id = %d ", params.ChatId)
	}

	query := `
		SELECT
			id,
			message,
			user_id,
			chat_id,
			created_at
		FROM chat_messages
	` + filter + `
		ORDER BY created_at DESC
	` + limit

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var message repo.ChatMessage

		err := rows.Scan(
			&message.ID,
			&message.Message,
			&message.UserId,
			&message.ChatId,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		message.UserInfo, err = getUserInfo(pr.db, message.UserId)
		if err != nil {
			return nil, err
		}

		result.Messages = append(result.Messages, &message)
	}

	queryCount := `SELECT count(1) FROM chat_messages` + filter
	err = pr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
