package integration

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/pmii/user-service/internal/domain/entity"
)

// TestChangePassword_Success tests admin changing user password
func TestChangePassword_Success(t *testing.T) {
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
		UserName:     "Target User",
		UserEmail:    "target@test.com",
		UserPassword: mustHashPassword(t, "oldpass123"),
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

	// Change password
	payload := map[string]interface{}{
		"oldPassword": "oldpass123",
		"newPassword": "newpass123",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users/"+strconv.Itoa(targetUser.UserId)+"/change-password", payload, adminToken)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Password changed successfully" {
		t.Errorf("expected 'Password changed successfully', got %s", resp.Meta.Message)
	}
}

// TestChangePassword_WrongOldPassword tests changing password with wrong old password
func TestChangePassword_WrongOldPassword(t *testing.T) {
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
		UserName:     "Target User",
		UserEmail:    "target@test.com",
		UserPassword: mustHashPassword(t, "oldpass123"),
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

	// Try to change password with wrong old password
	payload := map[string]interface{}{
		"oldPassword": "wrongpassword",
		"newPassword": "newpass123",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users/"+strconv.Itoa(targetUser.UserId)+"/change-password", payload, adminToken)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
