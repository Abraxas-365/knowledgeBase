package analiticssrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/opd/internal/analitics"
)

type Service struct {
	repo analitics.Repository
}

func NewService(repo analitics.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) GetAllAnalitics(ctx context.Context, startDate *time.Time, endDate *time.Time) ([]analitics.Statistic, error) {
	var allAnalitics []analitics.Statistic

	interactions, err := s.repo.GetInteractions(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	allAnalitics = append(allAnalitics, *interactions)

	data, err := s.repo.GetMostConsultedData(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	allAnalitics = append(allAnalitics, *data)

	users, err := s.repo.GetTotalUsers(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	allAnalitics = append(allAnalitics, *users)

	return allAnalitics, nil

}
