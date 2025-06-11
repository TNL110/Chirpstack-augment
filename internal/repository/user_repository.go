package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"go-auth-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, user.Email, user.PasswordHash, user.FullName).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, full_name, tenant_id, application_id, device_profile_id, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.FullName, &user.TenantID, &user.ApplicationID, &user.DeviceProfileID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, full_name, tenant_id, application_id, device_profile_id, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.FullName, &user.TenantID, &user.ApplicationID, &user.DeviceProfileID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetAllUsers retrieves all users with pagination
func (r *UserRepository) GetAllUsers(page, pageSize int) ([]models.User, int, error) {
	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM users"
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	offset := (page - 1) * pageSize
	query := `
		SELECT id, email, password_hash, full_name, tenant_id, application_id, device_profile_id, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash,
			&user.FullName, &user.TenantID, &user.ApplicationID, &user.DeviceProfileID,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return users, total, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Build dynamic query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Add updated_at with NOW() function
	setParts = append(setParts, "updated_at = NOW()")

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SearchUsers searches users by email or full name
func (r *UserRepository) SearchUsers(query string, page, pageSize int) ([]models.User, int, error) {
	searchPattern := "%" + strings.ToLower(query) + "%"

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*) FROM users
		WHERE LOWER(email) LIKE $1 OR LOWER(full_name) LIKE $1`
	err := r.db.QueryRow(countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get users with pagination
	offset := (page - 1) * pageSize
	searchQuery := `
		SELECT id, email, password_hash, full_name, tenant_id, application_id, device_profile_id, created_at, updated_at
		FROM users
		WHERE LOWER(email) LIKE $1 OR LOWER(full_name) LIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(searchQuery, searchPattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash,
			&user.FullName, &user.TenantID, &user.ApplicationID, &user.DeviceProfileID,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return users, total, nil
}

// UpdateUserChirpStackData updates ChirpStack related fields for a user
func (r *UserRepository) UpdateUserChirpStackData(id string, tenantID, applicationID, deviceProfileID string) error {
	query := `
		UPDATE users
		SET tenant_id = $1, application_id = $2, device_profile_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	result, err := r.db.Exec(query, tenantID, applicationID, deviceProfileID, id)
	if err != nil {
		return fmt.Errorf("failed to update user ChirpStack data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
