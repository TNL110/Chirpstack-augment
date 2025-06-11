package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	FullName        string    `json:"full_name" db:"full_name"`
	TenantID        *string   `json:"tenant_id,omitempty" db:"tenant_id"`
	ApplicationID   *string   `json:"application_id,omitempty" db:"application_id"`
	DeviceProfileID *string   `json:"device_profile_id,omitempty" db:"device_profile_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// User management request/response models
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	FullName string `json:"full_name" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

type UserListResponse struct {
	Users      []User `json:"users"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type UserSearchRequest struct {
	Query    string `form:"q"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

// Public user info (without sensitive data)
type PublicUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
