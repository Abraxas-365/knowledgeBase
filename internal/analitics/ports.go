package analitics

import (
	"context"
	"time"
)

type Repository interface {
	// GetInteractions returns statistics about interactions
	GetInteractions(ctx context.Context, startDate *time.Time, endDate *time.Time) (*Statistic, error)

	// GetMostConsultedData returns statistics about most consulted data
	GetMostConsultedData(ctx context.Context, startDate *time.Time, endDate *time.Time) (*Statistic, error)

	// GetTotalUsers returns statistics about total users
	GetTotalUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (*Statistic, error)

	GetDailyUsers(ctx context.Context, startDate, endDate time.Time) ([]DailyStatistic, error)
	GetDailyInteractions(ctx context.Context, startDate, endDate time.Time) ([]DailyStatistic, error)
	GetDailyActiveUsers(ctx context.Context, startDate, endDate time.Time, activeDays int) ([]DailyStatistic, error)
}
