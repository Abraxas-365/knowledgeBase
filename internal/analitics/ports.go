package analitics

import (
	"context"
	"time"

	"github.com/Abraxas-365/opd/internal/chatuser"
	"github.com/Abraxas-365/opd/internal/interaction"
	"github.com/Abraxas-365/opd/internal/kb"
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

	GetAllChatUsers(ctx context.Context, startDate, endDate *time.Time) ([]chatuser.ChatUser, error)
	GetAllInteractionsData(ctx context.Context, startDate, endDate *time.Time) ([]interaction.Interaction, error)
	GetAllFiles(ctx context.Context, startDate, endDate *time.Time) ([]kb.DataFile, error)
}
