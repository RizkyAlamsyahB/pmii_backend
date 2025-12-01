package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pmii/user-service/internal/domain/entity"
)

// TestCreateUser_Success tests admin creating a new user
func TestCreateUser_Success(t *testing.T) {
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

	// Create new user via API
	payload := map[string]interface{}{
		"fullName": "New User",
		"email":    "newuser@test.com",
		"password": "pass123",
		"level":    "2",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users", payload, adminToken)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Status != "success" {
		t.Errorf("expected success status, got %s", resp.Meta.Status)
	}

	var userData map[string]interface{}
	if err := json.Unmarshal(resp.Data, &userData); err != nil {
		t.Fatalf("failed to unmarshal user data: %v", err)
	}

	if userData["email"] != "newuser@test.com" {
		t.Errorf("expected email newuser@test.com, got %v", userData["email"])
	}
	if userData["role"] != "author" {
		t.Errorf("expected role author, got %v", userData["role"])
	}
}

// TestCreateUser_EmailAlreadyExists tests creating user with duplicate email
func TestCreateUser_EmailAlreadyExists(t *testing.T) {
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

	// Create existing user
	existingUser := &entity.User{
		UserName:     "Existing User",
		UserEmail:    "existing@test.com",
		UserPassword: mustHashPassword(t, "pass123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, existingUser); err != nil {
		t.Fatalf("failed to create existing user: %v", err)
	}

	// Get admin token
	adminToken, _ := app.jwtService.GenerateToken(
		adminUser.UserId,
		adminUser.UserEmail,
		entity.UserLevelAdmin,
	)

	// Try to create user with duplicate email
	payload := map[string]interface{}{
		"fullName": "Duplicate User",
		"email":    "existing@test.com",
		"password": "pass123",
		"level":    "2",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users", payload, adminToken)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestCreateUser_ValidationError tests creating user with invalid data
func TestCreateUser_ValidationError(t *testing.T) {
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

	// Try to create user with invalid email
	payload := map[string]interface{}{
		"fullName": "Test User",
		"email":    "invalid-email",
		"password": "pass123",
		"level":    "2",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users", payload, adminToken)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
