package user

import (
	"context"
	"errors"
	"time"

	"github.com/pmii/pmii-backend/shared/jwt"
	"github.com/pmii/user-service/internal/domain/entity"
	"github.com/pmii/user-service/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type user_usecase_impl struct {
	user_repo   repository.UserRepository
	jwt_service *jwt.JWTService
}

// NewUserUsecase creates new user usecase
func NewUserUsecase(userRepo repository.UserRepository, jwtSecret string, jwtExpireHours int) UserUsecase {
	jwtService := jwt.NewJWTService(jwtSecret, jwtExpireHours)
	
	return &user_usecase_impl{
		user_repo:   userRepo,
		jwt_service: jwtService,
	}
}

func (uc *user_usecase_impl) Register(ctx context.Context, name, email, password, level string) (*entity.User, error) {
	// Check if email already exists
	existingUser, _ := uc.user_repo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// Create user
	user := &entity.User{
		UserName:     name,
		UserEmail:    email,
		UserPassword: string(hashedPassword),
		UserLevel:    level,
		UserStatus:   entity.UserStatusActive,
		UserPhoto:    "",
	}
	
	err = uc.user_repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (uc *user_usecase_impl) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	// Find user by email
	user, err := uc.user_repo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}
	
	// Check if user is active
	if !user.IsActive() {
		return "", nil, errors.New("user account is inactive")
	}
	
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}
	
	// Generate JWT token
	token, err := uc.jwt_service.GenerateToken(user.UserId, user.UserEmail, user.UserLevel)
	if err != nil {
		return "", nil, err
	}
	
	return token, user, nil
}

func (uc *user_usecase_impl) GetById(ctx context.Context, userId int) (*entity.User, error) {
	return uc.user_repo.FindById(ctx, userId)
}

func (uc *user_usecase_impl) GetAll(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	offset := (page - 1) * limit
	
	users, err := uc.user_repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := uc.user_repo.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

func (uc *user_usecase_impl) Update(ctx context.Context, userId int, name, email, level, status, photo string) error {
	user, err := uc.user_repo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	
	user.UserName = name
	user.UserEmail = email
	user.UserLevel = level
	user.UserStatus = status
	if photo != "" {
		user.UserPhoto = photo
	}
	user.UpdatedAt = time.Now()
	
	return uc.user_repo.Update(ctx, user)
}

func (uc *user_usecase_impl) ChangePassword(ctx context.Context, userId int, oldPassword, newPassword string) error {
	user, err := uc.user_repo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	
	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}
	
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	user.UserPassword = string(hashedPassword)
	user.UpdatedAt = time.Now()
	
	return uc.user_repo.Update(ctx, user)
}

func (uc *user_usecase_impl) Delete(ctx context.Context, userId int) error {
	return uc.user_repo.Delete(ctx, userId)
}

func (uc *user_usecase_impl) GetProfile(ctx context.Context, userId int) (*entity.User, error) {
	return uc.user_repo.FindById(ctx, userId)
}

// generateToken generates JWT token (placeholder - will use shared/jwt)
func (uc *user_usecase_impl) generateToken(user *entity.User) string {
	// This will be implemented using shared/jwt package
	return "temporary-token-" + user.UserEmail
}
