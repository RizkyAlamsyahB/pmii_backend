package integration

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/pmii/user-service/internal/domain/entity"
)

// TestGetAllUsers_Success tests admin listing all users
func TestGetAllUsers_Success(t *testing.T) {
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

	// Create some test users
	for i := 1; i <= 3; i++ {
		user := &entity.User{
			UserName:     "Test User " + strconv.Itoa(i),
			UserEmail:    "user" + strconv.Itoa(i) + "@test.com",
			UserPassword: mustHashPassword(t, "pass123"),
			UserLevel:    entity.UserLevelAuthor,
			UserStatus:   entity.UserStatusActive,
		}
		if err := app.repo.Create(nil, user); err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}
	}

	// Get admin token
	adminToken, _ := app.jwtService.GenerateToken(
		adminUser.UserId,
		adminUser.UserEmail,
		entity.UserLevelAdmin,
	)

	// Get all users
	rr := app.performJSONRequest(t, http.MethodGet, "/v1/users?page=1&limit=10", nil, adminToken)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Status != "success" {
		t.Errorf("expected success status, got %s", resp.Meta.Status)
	}
}

// TestGetUserById_Success tests admin getting user by ID
func TestGetUserById_Success(t *testing.T) {
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

	// Get user by ID
	rr := app.performJSONRequest(t, http.MethodGet, "/v1/users/"+strconv.Itoa(targetUser.UserId), nil, adminToken)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	resp := decodeAPIResponse(t, rr)
	var userData map[string]interface{}
	if err := json.Unmarshal(resp.Data, &userData); err != nil {
		t.Fatalf("failed to unmarshal user data: %v", err)
	}

	if userData["email"] != "target@test.com" {
		t.Errorf("expected email target@test.com, got %v", userData["email"])
	}
}

// TestGetUserById_NotFound tests getting non-existent user
func TestGetUserById_NotFound(t *testing.T) {
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

	// Get non-existent user
	rr := app.performJSONRequest(t, http.MethodGet, "/v1/users/9999", nil, adminToken)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
