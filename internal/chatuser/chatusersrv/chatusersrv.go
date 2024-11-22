package chatusersrv

import (
	"context"

	"github.com/Abraxas-365/opd/internal/chatuser"
	"github.com/google/uuid"
)

type Service struct {
	repo chatuser.Repository
}

func New(repo chatuser.Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) CreateChatUser(ctx context.Context, cu chatuser.ChatUser) (*chatuser.ChatUser, error) {
	if cu.ID == nil {
		id := uuid.New().String()
		cu.ID = &id
	}
	return s.repo.CreateChatUser(ctx, cu)
}

func (s *Service) GetChatUserByID(ctx context.Context, chatUserID string) (*chatuser.ChatUser, error) {
	return s.repo.GetChatUserByID(ctx, chatUserID)
}
