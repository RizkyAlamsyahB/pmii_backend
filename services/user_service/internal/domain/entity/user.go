package entity

import "time"

// User represents user entity from tbl_user
type User struct {
	UserId       int       `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name" db:"user_name"`
	UserEmail    string    `json:"user_email" db:"user_email"`
	UserPassword string    `json:"-" db:"user_password"` // Hidden from JSON
	UserLevel    string    `json:"user_level" db:"user_level"` // 1=Admin, 2=Author
	UserStatus   string    `json:"user_status" db:"user_status"` // 1=Active, 0=Inactive
	UserPhoto    string    `json:"user_photo" db:"user_photo"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserLevel constants
const (
	UserLevelAdmin  = "1"
	UserLevelAuthor = "2"
)

// UserStatus constants
const (
	UserStatusActive   = "1"
	UserStatusInactive = "0"
)

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.UserLevel == UserLevelAdmin
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.UserStatus == UserStatusActive
}
