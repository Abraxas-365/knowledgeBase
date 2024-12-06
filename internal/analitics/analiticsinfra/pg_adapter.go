package analyticsinfra

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *PostgresStore) getDefaultDateRange() (start, end time.Time) {
	now := time.Now()
	start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	return start, now
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func (s *PostgresStore) GetInteractions(ctx context.Context, startDate *time.Time, endDate *time.Time) (*analitics.Statistic, error) {
	start, end := s.getDefaultDateRange()
	if startDate != nil && endDate != nil {
		start = *startDate
		end = *endDate
	}

	query := `
        SELECT COUNT(*) 
        FROM interactions 
        WHERE created_at >= $1 
        AND created_at <= $2`

	var count int64
	err := s.db.GetContext(ctx, &count, query, start, end)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("failed to get interactions count: %v", err))
	}

	if count == 0 {
		return nil, errors.ErrNotFound("no interactions found for the specified period")
	}

	if startDate == nil && endDate == nil {
		return analitics.NewTotalMonthlyInteracions(strconv.FormatInt(count, 10)), nil
	}

	return analitics.NewInteractionsBetweenDates(
		strconv.FormatInt(count, 10),
		formatDate(start),
		formatDate(end),
	), nil
}

func (s *PostgresStore) GetMostConsultedData(ctx context.Context, startDate *time.Time, endDate *time.Time) (*analitics.Statistic, error) {
	start, end := s.getDefaultDateRange()
	if startDate != nil && endDate != nil {
		start = *startDate
		end = *endDate
	}

	query := `
        SELECT filename 
        FROM files 
        WHERE created_at >= $1 
        AND created_at <= $2 
        GROUP BY filename 
        ORDER BY COUNT(*) DESC 
        LIMIT 1`

	var filename string
	err := s.db.GetContext(ctx, &filename, query, start, end)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound("no data consulted in the specified period")
		}
		return nil, errors.ErrDatabase(fmt.Sprintf("failed to get most consulted data: %v", err))
	}

	if startDate == nil && endDate == nil {
		return analitics.NewMonthlyMostConsultedData(filename), nil
	}

	return analitics.NewMostConsultedDataBetweenDates(
		filename,
		formatDate(start),
		formatDate(end),
	), nil
}

func (s *PostgresStore) GetTotalUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (*analitics.Statistic, error) {
	start, end := s.getDefaultDateRange()
	if startDate != nil && endDate != nil {
		start = *startDate
		end = *endDate
	}

	query := `
        SELECT COUNT(DISTINCT user_id) 
        FROM chatUser 
        WHERE created_at >= $1 
        AND created_at <= $2`

	var count int64
	err := s.db.GetContext(ctx, &count, query, start, end)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("failed to get total users: %v", err))
	}

	if count == 0 {
		return nil, errors.ErrNotFound("no users found for the specified period")
	}

	if startDate == nil && endDate == nil {
		return analitics.NewTotalUsers(strconv.FormatInt(count, 10)), nil
	}

	return analitics.NewUsersBetweenDates(
		strconv.FormatInt(count, 10),
		formatDate(start),
		formatDate(end),
	), nil
}

func (s *PostgresStore) GetDailyUsers(ctx context.Context, startDate, endDate time.Time) ([]analitics.DailyStatistic, error) {
	query := `
        WITH RECURSIVE dates AS (
            SELECT DATE(:start_date) AS date
            UNION ALL
            SELECT date + INTERVAL '1 day'
            FROM dates
            WHERE date < DATE(:end_date)
        )
        SELECT 
            dates.date::timestamp with time zone as date,
            COALESCE(COUNT(DISTINCT i.user_chat_id), 0) as count
        FROM dates
        LEFT JOIN interactions i ON DATE(i.created_at) = dates.date
        GROUP BY dates.date
        ORDER BY dates.date;
    `

	var stats []analitics.DailyStatistic
	args := map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	}

	err := s.db.SelectContext(ctx, &stats, s.db.Rebind(query), args)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("failed to get daily active users: %v", err))
	}

	return stats, nil
}

func (s *PostgresStore) GetDailyActiveUsers(ctx context.Context, startDate, endDate time.Time, activeDays int) ([]analitics.DailyStatistic, error) {
	query := `
        WITH RECURSIVE dates AS (
            SELECT DATE(:start_date) AS date
            UNION ALL
            SELECT date + INTERVAL '1 day'
            FROM dates
            WHERE date < DATE(:end_date)
        )
        SELECT 
            dates.date::timestamp with time zone as date,
            COALESCE(COUNT(DISTINCT user_chat_id), 0) as count
        FROM dates
        LEFT JOIN LATERAL (
            SELECT DISTINCT i.user_chat_id
            FROM interactions i
            WHERE i.created_at >= dates.date - INTERVAL '1 day' * :active_days
            AND i.created_at < dates.date + INTERVAL '1 day'
        ) active_users ON true
        GROUP BY dates.date
        ORDER BY dates.date;
    `

	var stats []analitics.DailyStatistic
	args := map[string]interface{}{
		"start_date":  startDate,
		"end_date":    endDate,
		"active_days": activeDays,
	}

	err := s.db.SelectContext(ctx, &stats, s.db.Rebind(query), args)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("failed to get daily active users: %v", err))
	}

	return stats, nil
}

// New function to get daily interactions
func (s *PostgresStore) GetDailyInteractions(ctx context.Context, startDate, endDate time.Time) ([]analitics.DailyStatistic, error) {
	query := `
        WITH RECURSIVE dates AS (
            SELECT DATE(:start_date) AS date
            UNION ALL
            SELECT date + INTERVAL '1 day'
            FROM dates
            WHERE date < DATE(:end_date)
        )
        SELECT 
            dates.date::timestamp with time zone as date,
            COALESCE(COUNT(i.id), 0) as count
        FROM dates
        LEFT JOIN interactions i ON DATE(i.created_at) = dates.date
        GROUP BY dates.date
        ORDER BY dates.date;
    `

	var stats []analitics.DailyStatistic
	args := map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	}

	err := s.db.SelectContext(ctx, &stats, s.db.Rebind(query), args)
	if err != nil {
		return nil, errors.ErrDatabase(fmt.Sprintf("failed to get daily interactions: %v", err))
	}

	return stats, nil
}
