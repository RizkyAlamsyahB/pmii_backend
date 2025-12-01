package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharedResponse "github.com/pmii/pmii-backend/shared/utils"
	"github.com/pmii/user-service/internal/delivery/http/dto/request"
	"github.com/pmii/user-service/internal/delivery/http/dto/response"
	"github.com/pmii/user-service/internal/domain/entity"
	"github.com/pmii/user-service/internal/usecase/user"
)

type UserHandler struct {
	user_usecase user.UserUsecase
}

// NewUserHandler creates new user handler
func NewUserHandler(userUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		user_usecase: userUsecase,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := map[string][]string{
			"validation": {err.Error()},
		}
		sharedResponse.ValidationErrorResponse(c, errors)
		return
	}

	user, err := h.user_usecase.Register(c.Request.Context(), req.UserName, req.UserEmail, req.Password, req.UserLevel)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	userResp := h.toUserResponse(user)
	sharedResponse.SuccessResponse(c, http.StatusCreated, "User registered successfully", userResp)
}

// Create handles admin creating new user
func (h *UserHandler) Create(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := map[string][]string{
			"validation": {err.Error()},
		}
		sharedResponse.ValidationErrorResponse(c, errors)
		return
	}

	user, err := h.user_usecase.Register(c.Request.Context(), req.UserName, req.UserEmail, req.Password, req.UserLevel)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	userResp := h.toUserResponse(user)
	sharedResponse.SuccessResponse(c, http.StatusCreated, "User created successfully", userResp)
}

// Login handles user authentication
func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := map[string][]string{
			"validation": {err.Error()},
		}
		sharedResponse.ValidationErrorResponse(c, errors)
		return
	}

	token, user, err := h.user_usecase.Login(c.Request.Context(), req.UserEmail, req.Password)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	loginResp := &response.LoginResponse{
		Token: token,
		User:  h.toUserResponse(user),
	}

	sharedResponse.SuccessResponse(c, http.StatusOK, "Login successful", loginResp)
}

// GetById retrieves user by ID
func (h *UserHandler) GetById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.user_usecase.GetById(c.Request.Context(), userId)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	userResp := h.toUserResponse(user)
	sharedResponse.SuccessResponse(c, http.StatusOK, "User retrieved successfully", userResp)
}

// GetAll retrieves all users with pagination
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.user_usecase.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	usersResp := make([]*response.UserResponse, 0)
	for _, user := range users {
		usersResp = append(usersResp, h.toUserResponse(user))
	}

	totalPages := (total + limit - 1) / limit
	paginationResp := &response.PaginationResponse{
		Meta: response.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Data: usersResp,
	}

	sharedResponse.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", paginationResp)
}

// Update updates user data
func (h *UserHandler) Update(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := map[string][]string{
			"validation": {err.Error()},
		}
		sharedResponse.ValidationErrorResponse(c, errors)
		return
	}

	err = h.user_usecase.Update(c.Request.Context(), userId, req.UserName, req.UserEmail, req.UserLevel, req.UserStatus, req.UserPhoto)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Failed to update user", err)
		return
	}

	sharedResponse.SuccessResponse(c, http.StatusOK, "User updated successfully", "")
}

// ChangePassword changes user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := map[string][]string{
			"validation": {err.Error()},
		}
		sharedResponse.ValidationErrorResponse(c, errors)
		return
	}

	err = h.user_usecase.ChangePassword(c.Request.Context(), userId, req.OldPassword, req.NewPassword)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Failed to change password", err)
		return
	}

	sharedResponse.SuccessResponse(c, http.StatusOK, "Password changed successfully", "")
}

// Delete deletes user
func (h *UserHandler) Delete(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	err = h.user_usecase.Delete(c.Request.Context(), userId)
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusBadRequest, "Failed to delete user", err)
		return
	}

	sharedResponse.SuccessResponse(c, http.StatusOK, "User deleted successfully", "")
}

// GetProfile retrieves current user profile (from JWT token)
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user_id from context (set by auth middleware)
	userId, exists := c.Get("user_id")
	if !exists {
		sharedResponse.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := h.user_usecase.GetProfile(c.Request.Context(), userId.(int))
	if err != nil {
		sharedResponse.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	userResp := h.toUserResponse(user)
	sharedResponse.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", userResp)
}

// Logout handles user logout (client-side token invalidation)
func (h *UserHandler) Logout(c *gin.Context) {
	// For JWT, logout is handled client-side by removing the token
	// Server can optionally implement token blacklist here
	sharedResponse.SuccessResponse(c, http.StatusOK, "Logout successful", "")
}

// toUserResponse converts entity.User to response.UserResponse
func (h *UserHandler) toUserResponse(user *entity.User) *response.UserResponse {
	// Convert user_level to role name
	role := "author"
	if user.UserLevel == entity.UserLevelAdmin {
		role = "admin"
	}

	// Convert user_status to status name
	status := "inactive"
	if user.UserStatus == entity.UserStatusActive {
		status = "active"
	}

	return &response.UserResponse{
		Id:        user.UserId,
		FullName:  user.UserName,
		Email:     user.UserEmail,
		Role:      role,
		Status:    status,
		Photo:     user.UserPhoto,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
