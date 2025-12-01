package integration

import (
	"net/http"
	"testing"

	"github.com/pmii/user-service/internal/domain/entity"
)

// TestCreateUser_Forbidden tests author cannot create user
func TestCreateUser_Forbidden(t *testing.T) {
	app := newTestApp(t)

	// Create author user
	authorUser := &entity.User{
		UserName:     "Author User",
		UserEmail:    "author@test.com",
		UserPassword: mustHashPassword(t, "author123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, authorUser); err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	// Get author token
	authorToken, _ := app.jwtService.GenerateToken(
		authorUser.UserId,
		authorUser.UserEmail,
		entity.UserLevelAuthor,
	)

	// Try to create user (should fail)
	payload := map[string]interface{}{
		"fullName": "Forbidden User",
		"email":    "forbidden@test.com",
		"password": "pass123",
		"level":    "2",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users", payload, authorToken)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Admin access only" {
		t.Errorf("expected 'Admin access only' message, got %s", resp.Meta.Message)
	}
}

// TestGetAllUsers_Forbidden tests author cannot list users
func TestGetAllUsers_Forbidden(t *testing.T) {
	app := newTestApp(t)

	// Create author user
	authorUser := &entity.User{
		UserName:     "Author User",
		UserEmail:    "author@test.com",
		UserPassword: mustHashPassword(t, "author123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, authorUser); err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	// Get author token
	authorToken, _ := app.jwtService.GenerateToken(
		authorUser.UserId,
		authorUser.UserEmail,
		entity.UserLevelAuthor,
	)

	// Try to get all users (should fail)
	rr := app.performJSONRequest(t, http.MethodGet, "/v1/users", nil, authorToken)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Admin access only" {
		t.Errorf("expected 'Admin access only', got %s", resp.Meta.Message)
	}
}

// TestUpdateUser_Forbidden tests author cannot update user
func TestUpdateUser_Forbidden(t *testing.T) {
	app := newTestApp(t)

	// Create author user
	authorUser := &entity.User{
		UserName:     "Author User",
		UserEmail:    "author@test.com",
		UserPassword: mustHashPassword(t, "author123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, authorUser); err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	// Get author token
	authorToken, _ := app.jwtService.GenerateToken(
		authorUser.UserId,
		authorUser.UserEmail,
		entity.UserLevelAuthor,
	)

	// Try to update user (should fail)
	payload := map[string]interface{}{
		"fullName": "Updated Name",
		"email":    "updated@test.com",
		"level":    "2",
		"status":   "1",
	}

	rr := app.performJSONRequest(t, http.MethodPut, "/v1/users/1", payload, authorToken)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Admin access only" {
		t.Errorf("expected 'Admin access only', got %s", resp.Meta.Message)
	}
}

// TestDeleteUser_Forbidden tests author cannot delete user
func TestDeleteUser_Forbidden(t *testing.T) {
	app := newTestApp(t)

	// Create author user
	authorUser := &entity.User{
		UserName:     "Author User",
		UserEmail:    "author@test.com",
		UserPassword: mustHashPassword(t, "author123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, authorUser); err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	// Get author token
	authorToken, _ := app.jwtService.GenerateToken(
		authorUser.UserId,
		authorUser.UserEmail,
		entity.UserLevelAuthor,
	)

	// Try to delete user (should fail)
	rr := app.performJSONRequest(t, http.MethodDelete, "/v1/users/1", nil, authorToken)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Admin access only" {
		t.Errorf("expected 'Admin access only', got %s", resp.Meta.Message)
	}
}

// TestChangePassword_Forbidden tests author cannot change other user's password
func TestChangePassword_Forbidden(t *testing.T) {
	app := newTestApp(t)

	// Create author user
	authorUser := &entity.User{
		UserName:     "Author User",
		UserEmail:    "author@test.com",
		UserPassword: mustHashPassword(t, "author123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(nil, authorUser); err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	// Get author token
	authorToken, _ := app.jwtService.GenerateToken(
		authorUser.UserId,
		authorUser.UserEmail,
		entity.UserLevelAuthor,
	)

	// Try to change password (should fail)
	payload := map[string]interface{}{
		"oldPassword": "oldpass",
		"newPassword": "newpass",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/users/1/change-password", payload, authorToken)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Admin access only" {
		t.Errorf("expected 'Admin access only', got %s", resp.Meta.Message)
	}
}

// TestUnauthorized_NoToken tests accessing protected endpoints without token
func TestUnauthorized_NoToken(t *testing.T) {
	app := newTestApp(t)

	// Try to access without token
	rr := app.performJSONRequest(t, http.MethodGet, "/v1/users", nil, "")

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
