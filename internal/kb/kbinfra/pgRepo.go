package kbinfra

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/Abraxas-365/opd/internal/kb"
	"github.com/Abraxas-365/toolkit/pkg/database"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{
		db,
	}
}

func (lc *PostgresStore) GetKnowlegeBaseConfig() (*kb.KnowlegeBaseConfig, error) {
	id := os.Getenv("KB_ID")
	if id == "" {
		return nil, errors.ErrUnexpected("KB_ID is not set")
	}

	numberOfResults := os.Getenv("KB_NUMBER_OF_RESULTS")
	if numberOfResults == "" {
		return nil, errors.ErrUnexpected("KB_NUMBER_OF_RESULTS is not set")
	}

	numberOfResultsInt, err := strconv.Atoi(numberOfResults)
	if err != nil {
		return nil, errors.ErrUnexpected("KB_NUMBER_OF_RESULTS is not a number")
	}

	region := os.Getenv("KB_REGION")
	if region == "" {
		return nil, errors.ErrUnexpected("KB_REGION is not set")
	}

	modelId := os.Getenv("KB_MODEL_ID")
	if modelId == "" {
		return nil, errors.ErrUnexpected("KB_MODEL_ID is not set")
	}

	modelPrompt := os.Getenv("KB_MODEL_PROMPT")
	if modelPrompt == "" {
		return nil, errors.ErrUnexpected("KB_MODEL_PROMPT is not set")
	}

	s3DataSource := os.Getenv("KB_S3_DATA_SOURCE")
	if s3DataSource == "" {
		return nil, errors.ErrUnexpected("KB_S3_DATA_SOURCE is not set")
	}

	return &kb.KnowlegeBaseConfig{
		ID:              id,
		NumberOfResults: numberOfResultsInt,
		Region:          region,
		S3DataSurce:     s3DataSource,
		Model: kb.ModelInformation{
			ModelId: modelId,
			Prompt:  modelPrompt,
		}}, nil

}

func (lc *PostgresStore) SaveData(ctx context.Context, data kb.DataFile) (*kb.DataFile, error) {
	query := `
        INSERT INTO files (filename, s3_key, user_id)
        VALUES ($1, $2, $3)
        RETURNING id, filename, s3_key, user_id`

	var savedFile kb.DataFile

	err := lc.db.QueryRowxContext(
		ctx,
		query,
		data.Filename,
		data.S3Key,
		data.UserID,
	).StructScan(&savedFile)

	if err != nil {
		return nil, errors.ErrDatabase("failed to save file data: " + err.Error())
	}

	return &savedFile, nil
}

func (lc *PostgresStore) DeleteData(ctx context.Context, dataId int) (*kb.DataFile, error) {
	query := `
        DELETE FROM files 
        WHERE id = $1
        RETURNING id, filename, s3_key, user_id`

	var deletedFile kb.DataFile

	err := lc.db.QueryRowxContext(
		ctx,
		query,
		dataId,
	).StructScan(&deletedFile)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound("file not found")
		}
		return nil, errors.ErrDatabase("failed to delete file data: " + err.Error())
	}

	return &deletedFile, nil
}

func (lc *PostgresStore) getTotalCount(ctx context.Context, query string) (int, error) {
	var count int
	err := lc.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, errors.ErrDatabase(fmt.Sprintf("Failed to get total count: %v", err))
	}
	return count, nil
}

func (lc *PostgresStore) GetData(ctx context.Context, page, pageSize int) (database.PaginatedRecord[kb.DataFile], error) {
	offset := (page - 1) * pageSize

	query := `
        SELECT id, filename, s3_key, user_id 
        FROM files 
        ORDER BY id DESC
        LIMIT $1 OFFSET $2`

	var files []kb.DataFile
	err := lc.db.SelectContext(ctx, &files, query, pageSize, offset)
	if err != nil {
		return database.PaginatedRecord[kb.DataFile]{},
			errors.ErrDatabase(fmt.Sprintf("Failed to get files: %v", err))
	}

	total, err := lc.getTotalCount(ctx, `SELECT COUNT(*) FROM files`)
	if err != nil {
		return database.PaginatedRecord[kb.DataFile]{}, err
	}

	return database.PaginatedRecord[kb.DataFile]{
		Data:       files,
		PageNumber: page,
		PageSize:   pageSize,
		Total:      total,
	}, nil
}

func (lc *PostgresStore) GetDataById(ctx context.Context, id int) (*kb.DataFile, error) {
	query := `
        SELECT id, filename, s3_key, user_id 
        FROM files 
        WHERE id = $1`

	var file kb.DataFile
	err := lc.db.QueryRowxContext(ctx, query, id).StructScan(&file)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound("file not found")
		}
		return nil, errors.ErrDatabase(fmt.Sprintf("Failed to get file: %v", err))
	}

	return &file, nil
}
