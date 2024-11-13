package userinfra

import (
	"context"
	"fmt"

	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/toolkit/pkg/database"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db *sqlx.DB
}

func (s *PostgresStore) GetUserByID(ctx context.Context, userID string) (*user.User, error) {
	var u user.User

	query := `SELECT id, email, is_admin, provider, provider_id FROM "user" WHERE id = $1`
	err := s.db.GetContext(ctx, &u, query, userID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// NewUserStore creates a new PostgresStore for user repository
func NewUserStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// GetUserByProviderID retrieves a user by provider and provider ID
func (s *PostgresStore) GetUserByProviderID(ctx context.Context, provider, providerID string) (*user.User, error) {
	query := `SELECT id, email, is_admin, provider, provider_id FROM "user" WHERE provider = $1 AND provider_id = $2`
	var u user.User
	err := s.db.GetContext(ctx, &u, query, provider, providerID)
	if err != nil {
		return nil, errors.ErrNotFound("User not found")
	}
	return &u, nil
}

// CreateUser inserts a new user
func (s *PostgresStore) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	query := `INSERT INTO "user" (id, email, provider, provider_id, is_admin) VALUES ($1, $2, $3, $4, $5) RETURNING id, email, is_admin`
	err := s.db.QueryRowContext(ctx, query, u.ID, u.Email, u.Provider, u.ProviderID, u.IsAdmin).Scan(&u.ID, &u.Email, &u.IsAdmin)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("Failed to create user: %v", err))
	}
	return u, nil
}

// GetUsers retrieves a paginated list of all users
func (s *PostgresStore) GetUsers(ctx context.Context, page, pageSize int) (database.PaginatedRecord[user.User], error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, email, is_admin FROM "user" LIMIT $1 OFFSET $2`
	var users []user.User
	err := s.db.SelectContext(ctx, &users, query, pageSize, offset)
	if err != nil {
		return database.PaginatedRecord[user.User]{}, errors.ErrDatabase(fmt.Sprintf("Failed to get users: %v", err))
	}
	total, err := s.getTotalCount(ctx, `SELECT COUNT(*) FROM "user"`)
	if err != nil {
		return database.PaginatedRecord[user.User]{}, err
	}
	return database.PaginatedRecord[user.User]{
		Data:       users,
		PageNumber: page,
		PageSize:   pageSize,
		Total:      total,
	}, nil
}

// GetNotAdminUsers retrieves users without admin privileges
func (s *PostgresStore) GetNotAdminUsers(ctx context.Context, page, pageSize int) (database.PaginatedRecord[user.User], error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, email, is_admin FROM "user" WHERE is_admin = false LIMIT $1 OFFSET $2`
	var users []user.User
	err := s.db.SelectContext(ctx, &users, query, pageSize, offset)
	if err != nil {
		return database.PaginatedRecord[user.User]{}, errors.ErrDatabase(fmt.Sprintf("Failed to get not-admin users: %v", err))
	}
	total, err := s.getTotalCount(ctx, `SELECT COUNT(*) FROM "user" WHERE is_admin = false`)
	if err != nil {
		return database.PaginatedRecord[user.User]{}, err
	}

	return database.PaginatedRecord[user.User]{
		Data:       users,
		PageNumber: page,
		PageSize:   pageSize,
		Total:      total,
	}, nil
}

// GetUsersAdminRole retrieves users with admin privileges
func (s *PostgresStore) GetUsersAdminRole(ctx context.Context, page, pageSize int) (database.PaginatedRecord[user.User], error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, email, is_admin FROM "user" WHERE is_admin = true LIMIT $1 OFFSET $2`
	var users []user.User
	err := s.db.SelectContext(ctx, &users, query, pageSize, offset)
	if err != nil {
		return database.PaginatedRecord[user.User]{}, errors.ErrDatabase(fmt.Sprintf("Failed to get admin users: %v", err))
	}

	total, err := s.getTotalCount(ctx, `SELECT COUNT(*) FROM "user" WHERE is_admin = true`)
	if err != nil {
		return database.PaginatedRecord[user.User]{}, err
	}

	return database.PaginatedRecord[user.User]{
		Data:       users,
		PageNumber: page,
		PageSize:   pageSize,
		Total:      total,
	}, nil
}

// UpdateUser updates user information
func (s *PostgresStore) UpdateUser(ctx context.Context, u *user.User) (*user.User, error) {
	query := `UPDATE "user" SET email = $2, is_admin = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, email, is_admin`
	err := s.db.QueryRowContext(ctx, query, u.ID, u.Email, u.IsAdmin).Scan(&u.ID, &u.Email, &u.IsAdmin)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("Failed to update user: %v", err))
	}
	return u, nil
}

// DeleteUser deletes a user by ID
func (s *PostgresStore) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM "user" WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.ErrDatabase(fmt.Sprintf("Failed to delete user: %v", err))
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound("User not found")
	}
	return nil
}

// GetWhitelist retrieves all whitelisted emails
func (s *PostgresStore) GetWhitelist(ctx context.Context) ([]string, error) {
	query := `SELECT email FROM email_whitelist`
	var whitelist []string
	err := s.db.SelectContext(ctx, &whitelist, query)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("Failed to get whitelist: %v", err))
	}
	return whitelist, nil
}

// AddToWhitelist adds an email to the whitelist
func (s *PostgresStore) AddToWhitelist(ctx context.Context, email string) error {
	query := `INSERT INTO email_whitelist (email) VALUES ($1)`
	_, err := s.db.ExecContext(ctx, query, email)
	if err != nil {
		return errors.ErrDatabase(fmt.Sprintf("Failed to add to whitelist: %v", err))
	}
	return nil
}

// RemoveFromWhitelist removes an email from the whitelist
func (s *PostgresStore) RemoveFromWhitelist(ctx context.Context, email string) error {
	query := `DELETE FROM email_whitelist WHERE email = $1`
	result, err := s.db.ExecContext(ctx, query, email)
	if err != nil {
		return errors.ErrDatabase(fmt.Sprintf("Failed to remove from whitelist: %v", err))
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound("Email not found in whitelist")
	}
	return nil
}

// IsInWhite checks if an email is in the whitelist
func (s *PostgresStore) IsInWhitelist(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM email_whitelist WHERE email = $1)`
	var exists bool
	err := s.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, errors.ErrDatabase(fmt.Sprintf("Failed to check whitelist: %v", err))
	}
	return exists, nil
}

// PromoteUserToAdmin promotes a user to admin
func (s *PostgresStore) PromoteUserToAdmin(ctx context.Context, userID string) error {
	query := `UPDATE "user" SET is_admin = true WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.ErrDatabase(fmt.Sprintf("Failed to promote user to admin: %v", err))
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound("User not found")
	}
	return nil
}

func (s *PostgresStore) getTotalCount(ctx context.Context, query string, args ...interface{}) (int, error) {
	var total int
	err := s.db.GetContext(ctx, &total, query, args...)
	if err != nil {
		return 0, errors.ErrDatabase(fmt.Sprintf("Failed to get total count: %v", err))
	}
	return total, nil
}
