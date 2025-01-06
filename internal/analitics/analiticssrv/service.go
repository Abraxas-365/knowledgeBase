package analiticssrv

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"
	"time"

	"github.com/Abraxas-365/opd/internal/analitics"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/Abraxas-365/toolkit/pkg/s3client"
)

type Service struct {
	repo     analitics.Repository
	s3Client s3client.Client
}

func NewService(repo analitics.Repository, s3Client s3client.Client) *Service {
	return &Service{
		repo:     repo,
		s3Client: s3Client,
	}
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

func (s Service) ExportDatabaseToCSV(ctx context.Context, startDate, endDate *time.Time) (string, error) {
	var start, end time.Time
	if startDate != nil && endDate != nil {
		start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		end = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
	}

	chatUsers, err := s.repo.GetAllChatUsers(ctx, startDate, endDate)
	if err != nil {
		return "", err
	}

	interactions, err := s.repo.GetAllInteractionsData(ctx, startDate, endDate)
	if err != nil {
		return "", err
	}

	files, err := s.repo.GetAllFiles(ctx, startDate, endDate)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// Write Chat Users
	writer.Write([]string{"ID", "Age", "Gender", "Occupation", "Location"})
	for _, chatUser := range chatUsers {
		id := ""
		if chatUser.ID != nil {
			id = *chatUser.ID
		}
		writer.Write([]string{
			id,
			strconv.Itoa(chatUser.Age),
			chatUser.Gender,
			chatUser.Ocupation,
			chatUser.Location,
		})
	}

	writer.Write([]string{""})

	// Write Interactions
	writer.Write([]string{"ID", "User Chat ID", "Context Interaction Amount"})
	for _, interaction := range interactions {
		writer.Write([]string{
			strconv.Itoa(interaction.ID),
			interaction.UserChatID,
			strconv.Itoa(len(interaction.ContextInteraction)),
		})
	}

	writer.Write([]string{""})

	// Write Files
	writer.Write([]string{"ID", "Filename", "S3 Key", "User ID", "User Email"})
	for _, file := range files {
		writer.Write([]string{
			string(rune(file.ID)),
			file.Filename,
			file.S3Key,
			file.UserID,
			file.UserEmail,
		})
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return "", errors.ErrUnexpected("error writing CSV: " + err.Error())
	}

	// Generate filename
	filenamePrefix := "complete"
	if startDate != nil && endDate != nil {
		filenamePrefix = start.Format("2006-01-02") + "_to_" + end.Format("2006-01-02")
	}
	filename := "exports/database_export_" + filenamePrefix + "_" + time.Now().Format("150405") + ".csv"

	// Upload to S3
	err = s.s3Client.UploadCSV(ctx, filename, buffer.Bytes())
	if err != nil {
		return "", errors.ErrServiceUnavailable("failed to upload to S3: " + err.Error())
	}

	// Generate presigned URL
	presignedURL, err := s.s3Client.GeneratePresignedGetURL(filename, 24*time.Hour)
	if err != nil {
		return "", errors.ErrServiceUnavailable("failed to generate presigned URL: " + err.Error())
	}

	return presignedURL, nil
}
