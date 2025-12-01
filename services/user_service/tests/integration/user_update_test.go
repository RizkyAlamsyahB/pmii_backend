package integration

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/pmii/user-service/internal/domain/entity"
)

// TestUpdateUser_Success tests admin updating a user
func TestUpdateUser_Success(t *testing.T) {
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
		UserName:     "Old Name",
		UserEmail:    "oldmail@test.com",
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

	// Update user
	payload := map[string]interface{}{
		"fullName": "New Name",
		"email":    "newmail@test.com",
		"level":    "2",
		"status":   "1",
	}

	rr := app.performJSONRequest(t, http.MethodPut, "/v1/users/"+strconv.Itoa(targetUser.UserId), payload, adminToken)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Status != "success" {
		t.Errorf("expected success status, got %s", resp.Meta.Status)
	}
}

// TestUpdateUser_NotFound tests updating non-existent user
func TestUpdateUser_NotFound(t *testing.T) {
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

	// Try to update non-existent user
	payload := map[string]interface{}{
		"fullName": "New Name",
		"email":    "newmail@test.com",
		"level":    "2",
		"status":   "1",
	}

	rr := app.performJSONRequest(t, http.MethodPut, "/v1/users/9999", payload, adminToken)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
