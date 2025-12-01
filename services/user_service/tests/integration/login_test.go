package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pmii/user-service/internal/delivery/http/dto/response"
	"github.com/pmii/user-service/internal/domain/entity"
)

func TestLogin_Success(t *testing.T) {
	app := newTestApp(t)

	user := &entity.User{
		UserName:     "Auth User",
		UserEmail:    "auth@example.com",
		UserPassword: mustHashPassword(t, "secret123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(context.Background(), user); err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	payload := map[string]string{
		"email":    user.UserEmail,
		"password": "secret123",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/auth/login", payload, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Login successful" {
		t.Fatalf("unexpected response message: %s", resp.Meta.Message)
	}

	var loginResp response.LoginResponse
	if err := json.Unmarshal(resp.Data, &loginResp); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}

	if loginResp.Token == "" {
		t.Fatalf("expected token to be generated")
	}
	if loginResp.User.Email != user.UserEmail {
		t.Errorf("expected user email %s, got %s", user.UserEmail, loginResp.User.Email)
	}
	if loginResp.User.Role != "author" {
		t.Errorf("expected role author, got %s", loginResp.User.Role)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	app := newTestApp(t)

	user := &entity.User{
		UserName:     "Auth User",
		UserEmail:    "auth@example.com",
		UserPassword: mustHashPassword(t, "secret123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(context.Background(), user); err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	payload := map[string]string{
		"email":    user.UserEmail,
		"password": "wrong-password",
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/auth/login", payload, "")

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Invalid credentials" {
		t.Fatalf("expected invalid credentials message, got %s", resp.Meta.Message)
	}
}
