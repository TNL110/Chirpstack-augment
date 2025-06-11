package interfaces

import "go-auth-api/internal/models"

type UserServiceInterface interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	GetUserByID(id string) (*models.User, error)
	GetAllUsers(page, pageSize int) (*models.UserListResponse, error)
	UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id string) error
	SearchUsers(req *models.UserSearchRequest) (*models.UserListResponse, error)
}
