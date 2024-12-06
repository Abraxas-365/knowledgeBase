package analyticsinfra

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/Abraxas-365/opd/internal/analitics"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewAnalyticsStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) GetInteractions(ctx context.Context, startDate *time.Time, endDate *time.Time) (*analitics.Statistic, error) {
	query := `SELECT COUNT(*) FROM interactions`
	args := []interface{}{}

	if startDate != nil && endDate != nil {
		query = `SELECT COUNT(*) FROM interactions WHERE created_at BETWEEN $1 AND $2`
		args = append(args, startDate, endDate)
	}

	var count int
	err := s.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get interactions count: " + err.Error())
	}

	if startDate != nil && endDate != nil {
		return analitics.NewInteractionsBetweenDates(
			strconv.Itoa(count),
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
		), nil
	}

	return analitics.NewTotalMonthlyInteracions(strconv.Itoa(count)), nil
}

func (s *PostgresStore) GetMostConsultedData(ctx context.Context, startDate *time.Time, endDate *time.Time) (*analitics.Statistic, error) {
	query := `
		WITH file_counts AS (
			SELECT f.filename, COUNT(*) as access_count 
			FROM files f
			JOIN interactions i ON f.id::text = ANY(i.context_interaction)
	`
	args := []interface{}{}

	if startDate != nil && endDate != nil {
		query += ` WHERE i.created_at BETWEEN $1 AND $2`
		args = append(args, startDate, endDate)
	}

	query += `
			GROUP BY f.filename
		)
		SELECT filename 
		FROM file_counts 
		ORDER BY access_count DESC 
		LIMIT 1
	`

	var filename string
	err := s.db.GetContext(ctx, &filename, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			filename = "No data available"
		} else {
			return nil, errors.ErrDatabase("failed to get most consulted data: " + err.Error())
		}
	}

	if startDate != nil && endDate != nil {
		return analitics.NewMostConsultedDataBetweenDates(
			filename,
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
		), nil
	}

	return analitics.NewMonthlyMostConsultedData(filename), nil
}

func (s *PostgresStore) GetTotalUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (*analitics.Statistic, error) {
	query := `SELECT COUNT(*) FROM "user"`
	args := []interface{}{}

	if startDate != nil && endDate != nil {
		query = `SELECT COUNT(*) FROM "user" WHERE created_at BETWEEN $1 AND $2`
		args = append(args, startDate, endDate)
	}

	var count int
	err := s.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get total users: " + err.Error())
	}

	if startDate != nil && endDate != nil {
		return analitics.NewUsersBetweenDates(
			strconv.Itoa(count),
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
		), nil
	}

	return analitics.NewTotalUsers(strconv.Itoa(count)), nil
}

func (s *PostgresStore) GetDailyUsers(ctx context.Context, startDate, endDate time.Time) ([]analitics.DailyStatistic, error) {
	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM "user"
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	var stats []analitics.DailyStatistic
	err := s.db.SelectContext(ctx, &stats, query, startDate, endDate)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get daily users: " + err.Error())
	}

	if len(stats) == 0 {
		return nil, errors.ErrNotFound("no daily user statistics found for the specified date range")
	}

	return stats, nil
}

func (s *PostgresStore) GetDailyActiveUsers(ctx context.Context, startDate, endDate time.Time, activeDays int) ([]analitics.DailyStatistic, error) {
	if activeDays < 1 {
		return nil, errors.ErrBadRequest("active days must be greater than 0")
	}

	query := `
		WITH daily_active AS (
			SELECT 
				DATE(created_at) as date,
				COUNT(DISTINCT user_chat_id) as count
			FROM interactions
			WHERE created_at BETWEEN $1 AND $2
			GROUP BY DATE(created_at)
			HAVING COUNT(DISTINCT user_chat_id) >= $3
		)
		SELECT date, count
		FROM daily_active
		ORDER BY date
	`

	var stats []analitics.DailyStatistic
	err := s.db.SelectContext(ctx, &stats, query, startDate, endDate, activeDays)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get daily active users: " + err.Error())
	}

	if len(stats) == 0 {
		return nil, errors.ErrNotFound("no daily active user statistics found for the specified criteria")
	}

	return stats, nil
}

func (s *PostgresStore) GetDailyInteractions(ctx context.Context, startDate, endDate time.Time) ([]analitics.DailyStatistic, error) {
	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM interactions
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	var stats []analitics.DailyStatistic
	err := s.db.SelectContext(ctx, &stats, query, startDate, endDate)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get daily interactions: " + err.Error())
	}

	if len(stats) == 0 {
		return nil, errors.ErrNotFound("no daily interaction statistics found for the specified date range")
	}

	return stats, nil
}
