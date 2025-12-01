package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pmii/user-service/internal/delivery/http/dto/response"
	"github.com/pmii/user-service/internal/domain/entity"
)

func TestRegister_Success(t *testing.T) {
	app := newTestApp(t)

	payload := map[string]string{
		"fullName": "John Doe",
		"email":    "john@example.com",
		"password": "secret123",
		"level":    entity.UserLevelAdmin,
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/auth/register", payload, "")

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "User registered successfully" {
		t.Fatalf("unexpected response message: %s", resp.Meta.Message)
	}

	var userResp response.UserResponse
	if err := json.Unmarshal(resp.Data, &userResp); err != nil {
		t.Fatalf("failed to decode user response: %v", err)
	}

	if userResp.Email != payload["email"] {
		t.Errorf("expected email %s, got %s", payload["email"], userResp.Email)
	}
	if userResp.FullName != payload["fullName"] {
		t.Errorf("expected fullName %s, got %s", payload["fullName"], userResp.FullName)
	}
	if userResp.Role != "admin" {
		t.Errorf("expected role admin, got %s", userResp.Role)
	}
	if userResp.Status != "active" {
		t.Errorf("expected status active, got %s", userResp.Status)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	app := newTestApp(t)

	existing := &entity.User{
		UserName:     "Existing",
		UserEmail:    "taken@example.com",
		UserPassword: mustHashPassword(t, "secret123"),
		UserLevel:    entity.UserLevelAuthor,
		UserStatus:   entity.UserStatusActive,
	}

	if err := app.repo.Create(context.Background(), existing); err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	payload := map[string]string{
		"fullName": "Jane",
		"email":    existing.UserEmail,
		"password": "secret123",
		"level":    entity.UserLevelAuthor,
	}

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/auth/register", payload, "")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Status != "error" {
		t.Fatalf("expected error status, got %s", resp.Meta.Status)
	}
	if resp.Meta.Message != "email already registered" {
		t.Fatalf("expected duplicate email message, got %s", resp.Meta.Message)
	}
}
