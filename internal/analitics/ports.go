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
}
