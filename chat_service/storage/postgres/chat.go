package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/telegram_clone/chat_service/pkg/utils"
	"gitlab.com/telegram_clone/chat_service/storage/repo"
)

type chatRepo struct {
	db *sqlx.DB
}

func NewChat(db *sqlx.DB) repo.ChatStrogeI {
	return &chatRepo{
		db: db,
	}
}

func (cr *chatRepo) Create(req *repo.CreateChatReq) (*repo.Chat, error) {
	chat := repo.Chat{
		Name:     req.Name,
		UserID:   req.UserID,
		ImageUrl: req.ImageUrl,
		ChatType: req.ChatType,
	}

	tx, err := cr.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	query := `
		INSERT INTO chats (
			name,
			user_id,
			chat_type,
			image_url
		) VALUES($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(
		query,
		req.Name,
		req.UserID,
		req.ChatType,
		utils.NullString(req.ImageUrl),
	).Scan(
		&chat.ID,
	)
	if err != nil {
		return nil, err
	}

	chat.UserInfo, err = getUserInfo(cr.db, req.UserID)
	if err != nil {
		return nil, err
	}

	queryMemeber := `
		INSERT INTO chat_members(
			chat_id,
			user_id
		) VALUES($1, $2)
	`

	for _, memberID := range req.Members {
		_, err = tx.Exec(queryMemeber, chat.ID, memberID)
		if err != nil {
			return nil, err
		}
	}

	return &chat, nil
}

func (pr *chatRepo) Get(id int64) (*repo.Chat, error) {
	var (
		chat     repo.Chat
		imageUrl sql.NullString
	)

	query := `
		SELECT
			id,
			name,
			user_id,
			chat_type,
			image_url
		FROM chats
		WHERE id=$1
	`

	err := pr.db.QueryRow(
		query,
		id,
	).Scan(
		&chat.ID,
		&chat.Name,
		&chat.UserID,
		&chat.ChatType,
		&imageUrl,
	)
	if err != nil {
		return nil, err
	}
	chat.ImageUrl = imageUrl.String
	chat.UserInfo, err = getUserInfo(pr.db, chat.UserID)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (pr *chatRepo) Update(chat *repo.Chat) (*repo.Chat, error) {
	query := `
		UPDATE chats SET
			name=$1,
			image_url=$2
		WHERE id=$3 and user_id = $4
		RETURNING 
			chat_type		
	`

	err := pr.db.QueryRow(
		query,
		chat.Name,
		utils.NullString(chat.ImageUrl),
		chat.ID,
		chat.UserID,
	).Scan(
		&chat.ChatType,
	)
	if err != nil {
		return nil, err
	}
	chat.UserInfo, err = getUserInfo(pr.db, chat.UserID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (pr *chatRepo) Delete(id, user_id int64) error {
	query := `DELETE FROM chats WHERE id=$1 and user_id = $2`

	result, err := pr.db.Exec(query, id, user_id)
	if err != nil {
		return err
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pr *chatRepo) GetAll(params *repo.GetAllChatsParams) (*repo.GetAllChats, error) {
	var imageUrl sql.NullString
	result := repo.GetAllChats{
		Chats: make([]*repo.Chat, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""
	if params.UserID > 0 {
		filter = fmt.Sprintf(" WHERE cm.user_id = %d ", params.UserID)
	}

	query := `
		SELECT 
			c.id,
			c.name ,
			c.user_id ,
			c.chat_type ,
			c.image_url
		FROM chats c INNER JOIN chat_members cm ON cm.chat_id = c.id
	` + filter + limit

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat repo.Chat

		err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.UserID,
			&chat.ChatType,
			&imageUrl,
		)
		if err != nil {
			return nil, err
		}

		chat.ImageUrl = imageUrl.String
		chat.UserInfo, err = getUserInfo(pr.db, chat.UserID)
		if err != nil {
			return nil, err
		}
		result.Chats = append(result.Chats, &chat)
	}

	queryCount := ` SELECT count(1) FROM chats c INNER JOIN chat_members cm ON cm.chat_id = c.id ` + filter
	err = pr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getUserInfo(db *sqlx.DB, id int64) (*repo.GetUserInfo, error) {
	var (
		username, profileImageUrl sql.NullString
	)
	query := `
		SELECT
			first_name,
			last_name,
			email,
			username,
			profile_image_url,
			created_at
		FROM users
		WHERE id=$1
	`
	var result repo.GetUserInfo
	err := db.QueryRow(
		query,
		id,
	).Scan(
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&username,
		&profileImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	result.ImageUrl = profileImageUrl.String
	result.UserName = username.String

	return &result, nil
}

func (chr *chatRepo) AddMember(req *repo.AddMemberRequest) error {
	query := `
		INSERT INTO chat_members(
			chat_id,
			user_id
		) VALUES($1,$2)
		`
	_, err := chr.db.Exec(query, req.ChatId, req.UserId)
	if err != nil {
		return err
	}

	return nil

}

func (chr *chatRepo) RemoveMember(req *repo.RemoveMemberRequest) error {
	query := `DELETE FROM chat_members where user_id=$1 and chat_id=$2`
	effect, err := chr.db.Exec(query, req.UserId, req.ChatId)
	if err != nil {
		return err
	}

	if count, _ := effect.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (ur *chatRepo) GetChatMembers(params *repo.GetChatMembersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	query := `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.username,
			u.profile_image_url,
			u.created_at
		FROM users u
		INNER JOIN chat_members cm ON cm.user_id=u.id
		WHERE cm.chat_id=$1
		ORDER BY u.first_name ASC, u.last_name ASC
		` + limit

	rows, err := ur.db.Query(query, params.ChatID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			u                         repo.User
			username, profileImageUrl sql.NullString
		)

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&username,
			&profileImageUrl,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		u.Username = username.String
		u.ProfileImageUrl = profileImageUrl.String

		result.Users = append(result.Users, &u)
	}

	queryCount := `
		SELECT count(1) FROM users u
		INNER JOIN chat_members cm ON cm.user_id=u.id
		WHERE cm.chat_id=$1`
	err = ur.db.QueryRow(queryCount, params.ChatID).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
