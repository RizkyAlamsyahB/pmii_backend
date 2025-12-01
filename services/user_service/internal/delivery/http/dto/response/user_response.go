package response

import "time"

// UserResponse for user data response (camelCase untuk FE)
type UserResponse struct {
	Id        int       `json:"id"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // "admin" atau "author"
	Status    string    `json:"status"` // "active" atau "inactive"
	Photo     string    `json:"photo,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginResponse for authentication response
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// PaginationMeta for pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// PaginationResponse for paginated data
type PaginationResponse struct {
	Meta PaginationMeta `json:"meta"`
	Data interface{}    `json:"data"`
}
