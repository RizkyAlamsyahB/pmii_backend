package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmii/pmii-backend/shared/jwt"
	"github.com/pmii/pmii-backend/shared/logger"
	"github.com/pmii/user-service/internal/delivery/http/handler"
	routerpkg "github.com/pmii/user-service/internal/delivery/http/router"
	"github.com/pmii/user-service/internal/domain/entity"
	"github.com/pmii/user-service/internal/domain/repository"
	"github.com/pmii/user-service/internal/usecase/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	testJWTSecret     = "integration-secret"
	testJWTExpireHour = 24
)

func init() {
	gin.SetMode(gin.TestMode)
	logger.Log = zap.NewNop()
}

type testApp struct {
	router     *gin.Engine
	repo       *fakeUserRepository
	jwtService *jwt.JWTService
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()

	repo := newFakeUserRepository()
	jwtService := jwt.NewJWTService(testJWTSecret, testJWTExpireHour)
	userUsecase := user.NewUserUsecase(repo, testJWTSecret, testJWTExpireHour)
	userHandler := handler.NewUserHandler(userUsecase)

	return &testApp{
		router:     routerpkg.SetupRouter(userHandler, jwtService),
		repo:       repo,
		jwtService: jwtService,
	}
}

func (app *testApp) performJSONRequest(t *testing.T, method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()

	var body []byte
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		body = encoded
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	app.router.ServeHTTP(rr, req)
	return rr
}

type apiResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Data   json.RawMessage     `json:"data"`
	Errors map[string][]string `json:"errors"`
}

func decodeAPIResponse(t *testing.T, rr *httptest.ResponseRecorder) apiResponse {
	t.Helper()

	var resp apiResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return resp
}

func mustHashPassword(t *testing.T, plain string) string {
	t.Helper()

	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hashed)
}

var (
	errUserNotFound = errors.New("user not found")
)

type fakeUserRepository struct {
	mu     sync.RWMutex
	users  map[int]*entity.User
	nextID int
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users:  make(map[int]*entity.User),
		nextID: 1,
	}
}

func (r *fakeUserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cloned := cloneUser(user)
	cloned.UserId = r.nextID
	if cloned.CreatedAt.IsZero() {
		cloned.CreatedAt = now
	}
	cloned.UpdatedAt = cloned.CreatedAt

	r.users[cloned.UserId] = cloned
	r.nextID++

	user.UserId = cloned.UserId
	user.CreatedAt = cloned.CreatedAt
	user.UpdatedAt = cloned.UpdatedAt

	return nil
}

func (r *fakeUserRepository) FindById(ctx context.Context, userId int) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.users[userId]; ok {
		return cloneUser(user), nil
	}
	return nil, errUserNotFound
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.UserEmail == email {
			return cloneUser(user), nil
		}
	}
	return nil, errUserNotFound
}

func (r *fakeUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		result = append(result, cloneUser(user))
	}

	// Apply offset and limit boundaries
	if offset > len(result) {
		return []*entity.User{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (r *fakeUserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, ok := r.users[user.UserId]
	if !ok {
		return errUserNotFound
	}

	stored.UserName = user.UserName
	stored.UserEmail = user.UserEmail
	stored.UserLevel = user.UserLevel
	stored.UserStatus = user.UserStatus
	stored.UserPhoto = user.UserPhoto
	stored.UpdatedAt = time.Now()
	return nil
}

func (r *fakeUserRepository) Delete(ctx context.Context, userId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[userId]; !ok {
		return errUserNotFound
	}
	delete(r.users, userId)
	return nil
}

func (r *fakeUserRepository) CountAll(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users), nil
}

func cloneUser(user *entity.User) *entity.User {
	if user == nil {
		return nil
	}
	cloned := *user
	return &cloned
}

var _ repository.UserRepository = (*fakeUserRepository)(nil)
