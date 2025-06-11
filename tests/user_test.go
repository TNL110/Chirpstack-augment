package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-auth-api/internal/auth"
	"go-auth-api/internal/handlers"
	"go-auth-api/internal/interfaces"
	"go-auth-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserService
type MockUserService struct {
	mock.Mock
}

// Implement UserServiceInterface
var _ interfaces.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockUserService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(page, pageSize int) (*models.UserListResponse, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*models.UserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(id, req)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) SearchUsers(req *models.UserSearchRequest) (*models.UserListResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*models.UserListResponse), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockUserService) {
	gin.SetMode(gin.TestMode)

	mockUserService := &MockUserService{}
	userHandler := handlers.NewUserHandler(mockUserService)

	router := gin.New()
	api := router.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	return router, mockUserService
}

func TestUserRegister(t *testing.T) {
	router, mockService := setupTestRouter()

	// Test successful registration
	t.Run("Successful Registration", func(t *testing.T) {
		userID := uuid.New()
		expectedResponse := &models.AuthResponse{
			Token: "test-jwt-token",
			User: models.User{
				ID:       userID,
				Email:    "test@example.com",
				FullName: "Test User",
			},
		}

		mockService.On("Register", mock.AnythingOfType("*models.RegisterRequest")).Return(expectedResponse, nil)

		reqBody := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			FullName: "Test User",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-jwt-token", response.Token)
		assert.Equal(t, "test@example.com", response.User.Email)

		mockService.AssertExpectations(t)
	})

	// Test invalid request body
	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserLogin(t *testing.T) {
	router, mockService := setupTestRouter()

	// Test successful login
	t.Run("Successful Login", func(t *testing.T) {
		userID := uuid.New()
		expectedResponse := &models.AuthResponse{
			Token: "test-jwt-token",
			User: models.User{
				ID:       userID,
				Email:    "test@example.com",
				FullName: "Test User",
			},
		}

		mockService.On("Login", mock.AnythingOfType("*models.LoginRequest")).Return(expectedResponse, nil)

		reqBody := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-jwt-token", response.Token)

		mockService.AssertExpectations(t)
	})
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Test password hashing
	hashedPassword, err := auth.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Test password verification
	isValid := auth.CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid)

	// Test wrong password
	isValid = auth.CheckPasswordHash("wrongpassword", hashedPassword)
	assert.False(t, isValid)
}

func TestJWTService(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret-key")
	userID := uuid.New()
	email := "test@example.com"

	// Test token generation
	token, err := jwtService.GenerateToken(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test token validation
	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)

	// Test invalid token
	_, err = jwtService.ValidateToken("invalid-token")
	assert.Error(t, err)
}
