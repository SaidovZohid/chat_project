package repo

import "time"

const (
	UserTypeSuperadmin = "superadmin"
	UserTypeUser       = "user"
)

type User struct {
	ID              int64
	FirstName       string
	LastName        string
	Email           string
	Password        string
	Username        string
	ProfileImageUrl string
	Type            string
	CreatedAt       time.Time
}

type GetAllUsersParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetAllUsersResult struct {
	Users []*User
	Count int32
}

type UpdatePassword struct {
	UserID   int64
	Password string
}

type SetUserImageRequest struct {
	UserId   int64
	ImageUrl string
}

type UserStorageI interface {
	Create(u *User) (*User, error)
	Get(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll(params *GetAllUsersParams) (*GetAllUsersResult, error)
	UpdatePassword(req *UpdatePassword) error
	Update(u *User) (*User, error)
	Delete(id int64) error
	SetUserImage(*SetUserImageRequest) (*User, error)
}
