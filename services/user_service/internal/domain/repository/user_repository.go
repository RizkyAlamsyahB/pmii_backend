package repository

import (
	"context"
	"github.com/pmii/user-service/internal/domain/entity"
)

// UserRepository interface for user data access
type UserRepository interface {
	// Create new user
	Create(ctx context.Context, user *entity.User) error
	
	// FindByID finds user by user_id
	FindById(ctx context.Context, userId int) (*entity.User, error)
	
	// FindByEmail finds user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// FindAll retrieves all users with pagination
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
	
	// Update user data
	Update(ctx context.Context, user *entity.User) error
	
	// Delete user by user_id
	Delete(ctx context.Context, userId int) error
	
	// CountAll counts total users
	CountAll(ctx context.Context) (int, error)
}
