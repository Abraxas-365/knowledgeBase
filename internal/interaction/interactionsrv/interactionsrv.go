package interactionsrv

import (
	"context"

	"github.com/Abraxas-365/opd/internal/interaction"
)

type Service struct {
	repo interaction.Repository
}

func New(repo interaction.Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) CreateInteraction(ctx context.Context, cu interaction.Interaction) (*interaction.Interaction, error) {
	return s.repo.CreateInteraction(ctx, cu)
}
