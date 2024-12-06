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

	// data, err := s.repo.GetMostConsultedData(ctx, startDate, endDate)
	// if err != nil {
	// 	return nil, err
	// }
	// allAnalitics = append(allAnalitics, *data)

	users, err := s.repo.GetTotalUsers(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	allAnalitics = append(allAnalitics, *users)

	return allAnalitics, nil

}

func (s Service) GetDailyUsersInRange(ctx context.Context, startDate, endDate time.Time) ([]analitics.DailyStatistic, error) {
	// Normalize times to start and end of day
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)

	return s.repo.GetDailyUsers(ctx, start, end)
}

// GetDailyInteractionsInRange gets interactions per day for a specific date range
func (s Service) GetDailyInteractionsInRange(ctx context.Context, startDate, endDate time.Time) ([]analitics.DailyStatistic, error) {
	// Normalize times to start and end of day
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)

	return s.repo.GetDailyInteractions(ctx, start, end)
}
