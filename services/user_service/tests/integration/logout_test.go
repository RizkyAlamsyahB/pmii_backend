package integration

import (
	"net/http"
	"testing"
)

func TestLogout_Success(t *testing.T) {
	app := newTestApp(t)

	rr := app.performJSONRequest(t, http.MethodPost, "/v1/auth/logout", nil, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	resp := decodeAPIResponse(t, rr)
	if resp.Meta.Message != "Logout successful" {
		t.Fatalf("expected logout success message, got %s", resp.Meta.Message)
	}
	if string(resp.Data) != `""` {
		t.Fatalf("expected data to be empty string, got %s", string(resp.Data))
	}
}
