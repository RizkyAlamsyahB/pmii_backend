package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pmii/user-service/internal/delivery/http/dto/response"
	"github.com/pmii/user-service/internal/domain/entity"
)

func TestMe_Success(t *testing.T) {
	app := newTestApp(t)

	user := &entity.User{
		UserName:     "Profile User",
		UserEmail:    "profile@example.com",
		UserPassword: mustHashPassword(t, "secret123"),
		UserLevel:    entity.UserLevelAdmin,
		UserStatus:   entity.UserStatusActive,
	}
	if err := app.repo.Create(context.Background(), user); err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	token, err := app.jwtService.GenerateToken(user.UserId, user.UserEmail, user.UserLevel)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	rr := app.performJSONRequest(t, http.MethodGet, "/v1/auth/me", nil, token)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)

	var userResp response.UserResponse
	if err := json.Unmarshal(resp.Data, &userResp); err != nil {
		t.Fatalf("failed to decode user response: %v", err)
	}

	if userResp.Email != user.UserEmail {
		t.Errorf("expected email %s, got %s", user.UserEmail, userResp.Email)
	}
	if userResp.Role != "admin" {
		t.Errorf("expected role admin, got %s", userResp.Role)
	}
}

func TestMe_MissingAuthorization(t *testing.T) {
	app := newTestApp(t)

	rr := app.performJSONRequest(t, http.MethodGet, "/v1/auth/me", nil, "")

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Authorization header required" {
		t.Fatalf("expected authorization header message, got %s", resp.Meta.Message)
	}
}
