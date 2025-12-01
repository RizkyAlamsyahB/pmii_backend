package request

// RegisterRequest for user registration
type RegisterRequest struct {
	UserName  string `json:"fullName" binding:"required,min=3,max=100"`
	UserEmail string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	UserLevel string `json:"level" binding:"required,oneof=1 2"` // 1=Admin, 2=Author
}

// LoginRequest for user authentication
type LoginRequest struct {
	UserEmail string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

// UpdateUserRequest for updating user data
type UpdateUserRequest struct {
	UserName   string `json:"fullName" binding:"required,min=3,max=100"`
	UserEmail  string `json:"email" binding:"required,email"`
	UserLevel  string `json:"level" binding:"required,oneof=1 2"`
	UserStatus string `json:"status" binding:"required,oneof=0 1"`
	UserPhoto  string `json:"photo"`
}

// ChangePasswordRequest for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
