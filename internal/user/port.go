package user

import (
	"context"

	"github.com/Abraxas-365/toolkit/pkg/database"
)

type Repository interface {
	GetUserByProviderID(ctx context.Context, provider, providerID string) (*User, error)
	CreateUser(ctx context.Context, u *User) (*User, error)
	GetUsers(ctx context.Context, page, pageSize int) (database.PaginatedRecord[User], error)
	GetNotAdminUsers(ctx context.Context, page, pageSize int) (database.PaginatedRecord[User], error)
	GetUsersAdminRole(ctx context.Context, page, pageSize int) (database.PaginatedRecord[User], error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, userID string) error
	GetWhitelist(ctx context.Context) ([]string, error)
	AddToWhitelist(ctx context.Context, email string) error
	RemoveFromWhitelist(ctx context.Context, email string) error
	IsInWhitelist(ctx context.Context, email string) (bool, error)
	PromoteUserToAdmin(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (*User, error)
}
