package integration

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/pmii/user-service/internal/domain/entity"
)

// TestDeleteUser_Success tests admin deleting a user
func TestDeleteUser_Success(t *testing.T) {
	app := newTestApp(t)

	// Create admin user
	adminUser := &entity.User{
		UserName:     "Admin User",
		UserEmail:    "admin@test.com",
		UserPassword: mustHashPassword(t, "admin123"),
		UserLevel:    entity.UserLevelAdmin,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, adminUser); err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}

	// Create target user
	targetUser := &entity.User{
		UserName:     "To Delete",
		UserEmail:    "delete@test.com",
		UserPassword: mustHashPassword(t, "pass123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, targetUser); err != nil {
		t.Fatalf("failed to create target user: %v", err)
	}

	// Get admin token
	adminToken, _ := app.jwtService.GenerateToken(
		adminUser.UserId,
		adminUser.UserEmail,
		entity.UserLevelAdmin,
	)

	// Delete user
	rr := app.performJSONRequest(t, http.MethodDelete, "/v1/users/"+strconv.Itoa(targetUser.UserId), nil, adminToken)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "User deleted successfully" {
		t.Errorf("expected 'User deleted successfully', got %s", resp.Meta.Message)
	}

	// Verify user is deleted
	_, err := app.repo.FindById(nil, targetUser.UserId)
	if err != errUserNotFound {
		t.Errorf("expected user to be deleted, but still exists")
	}
}

// TestDeleteUser_NotFound tests deleting non-existent user
func TestDeleteUser_NotFound(t *testing.T) {
	app := newTestApp(t)

	// Create admin user
	adminUser := &entity.User{
		UserName:     "Admin User",
		UserEmail:    "admin@test.com",
		UserPassword: mustHashPassword(t, "admin123"),
		UserLevel:    entity.UserLevelAdmin,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, adminUser); err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}

	// Get admin token
	adminToken, _ := app.jwtService.GenerateToken(
		adminUser.UserId,
		adminUser.UserEmail,
		entity.UserLevelAdmin,
	)

	// Try to delete non-existent user
	rr := app.performJSONRequest(t, http.MethodDelete, "/v1/users/9999", nil, adminToken)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
