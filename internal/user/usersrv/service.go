package usersrv

import (
	"context"
	"log"

	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/toolkit/pkg/database"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/Abraxas-365/toolkit/pkg/lucia"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetUser(ctx context.Context, userID string) (*user.User, error) {
	log.Printf("Fetching user with ID: %s", userID)
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("User with ID %s not found", userID)
			return nil, errors.ErrNotFound("user not found")
		}
		log.Printf("Database error fetching user by ID %s: %v", userID, err)
	}

	log.Printf("Successfully fetched user: %+v", u)
	return u, nil
}

func (s *Service) GetUserByProviderID(ctx context.Context, provider, providerID string) (*user.User, error) {
	u, err := s.repo.GetUserByProviderID(ctx, provider, providerID)
	if err != nil {
		return nil, err
	}
	log.Printf("User: %v", u)
	return u, nil
}

func (s *Service) CreateUser(ctx context.Context, userInfo *lucia.UserInfo) (*user.User, error) {
	isInBlacklisted, err := s.repo.IsInBlacklist(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	} else if isInBlacklisted {
		return nil, errors.ErrUnauthorized("email is blacklisted")
	}

	u := &user.User{
		ID:         lucia.GenerateID(),
		Email:      userInfo.Email,
		IsAdmin:    false,
		Provider:   userInfo.Provider,
		ProviderID: userInfo.ID,
	}

	u, err = s.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) GetUsers(ctx context.Context, page, pageSize int) (database.PaginatedRecord[user.User], error) {
	return s.repo.GetUsers(ctx, page, pageSize)
}

func (s *Service) GetNotAdminUsers(ctx context.Context, page, pageSize int) (database.PaginatedRecord[user.User], error) {
	return s.repo.GetNotAdminUsers(ctx, page, pageSize)
}

func (s *Service) GetUsersAdminRole(ctx context.Context, page, pageSize int) (database.PaginatedRecord[user.User], error) {
	return s.repo.GetUsersAdminRole(ctx, page, pageSize)
}

func (s *Service) PromoteUserToAdmin(ctx context.Context, userID string) error {
	return s.repo.PromoteUserToAdmin(ctx, userID)
}

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *Service) GetBlacklist(ctx context.Context) ([]string, error) {
	return s.repo.GetBlacklist(ctx)
}

func (s *Service) AddToBlacklist(ctx context.Context, email string) error {
	return s.repo.AddToBlacklist(ctx, email)
}

func (s *Service) RemoveFromBlacklist(ctx context.Context, email string) error {
	return s.repo.RemoveFromBlacklist(ctx, email)
}
