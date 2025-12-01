package user

import (
	"context"
	"github.com/pmii/user-service/internal/domain/entity"
)

// UserUsecase interface for user business logic
type UserUsecase interface {
	// Register creates new user account
	Register(ctx context.Context, name, email, password, level string) (*entity.User, error)
	
	// Login authenticates user and returns token
	Login(ctx context.Context, email, password string) (token string, user *entity.User, err error)
	
	// GetById retrieves user by ID
	GetById(ctx context.Context, userId int) (*entity.User, error)
	
	// GetAll retrieves all users with pagination
	GetAll(ctx context.Context, page, limit int) ([]*entity.User, int, error)
	
	// Update updates user data
	Update(ctx context.Context, userId int, name, email, level, status, photo string) error
	
	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userId int, oldPassword, newPassword string) error
	
	// Delete deletes user
	Delete(ctx context.Context, userId int) error
	
	// GetProfile gets current user profile by ID
	GetProfile(ctx context.Context, userId int) (*entity.User, error)
}
