package service

import (
	"fmt"
	"math"

	"go-auth-api/internal/auth"
	"go-auth-api/internal/models"
	"go-auth-api/internal/repository"
)

type UserService struct {
	userRepo          *repository.UserRepository
	jwtService        *auth.JWTService
	chirpStackService *ChirpStackService
}

func NewUserService(userRepo *repository.UserRepository, jwtService *auth.JWTService, chirpStackService *ChirpStackService) *UserService {
	return &UserService{
		userRepo:          userRepo,
		jwtService:        jwtService,
		chirpStackService: chirpStackService,
	}
}

func (s *UserService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create ChirpStack resources if enabled
	if s.chirpStackService != nil && s.chirpStackService.IsEnabled() {
		chirpStackData, err := s.chirpStackService.CreateUserResources(user.ID.String(), user.Email)
		if err != nil {
			// Log error but don't fail user registration
			fmt.Printf("Warning: Failed to create ChirpStack resources for user %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("Successfully created ChirpStack resources for user %s: TenantID=%s, ApplicationID=%s, DeviceProfileID=%s\n",
				user.Email, chirpStackData.TenantID, chirpStackData.ApplicationID, chirpStackData.DeviceProfileID)
		}
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *UserService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(page, pageSize int) (*models.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.GetAllUsers(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Email != "" && req.Email != existingUser.Email {
		// Check if email already exists
		emailUser, _ := s.userRepo.GetUserByEmail(req.Email)
		if emailUser != nil && emailUser.ID.String() != id {
			return nil, fmt.Errorf("email already exists")
		}
		updates["email"] = req.Email
	}

	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}

	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password_hash"] = hashedPassword
	}

	if len(updates) == 0 {
		return existingUser, nil
	}

	// Update user
	err = s.userRepo.UpdateUser(id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Return updated user
	return s.userRepo.GetUserByID(id)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Delete user
	err = s.userRepo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// SearchUsers searches users by query with pagination
func (s *UserService) SearchUsers(req *models.UserSearchRequest) (*models.UserListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}
	if req.Query == "" {
		return s.GetAllUsers(req.Page, req.PageSize)
	}

	users, total, err := s.userRepo.SearchUsers(req.Query, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
