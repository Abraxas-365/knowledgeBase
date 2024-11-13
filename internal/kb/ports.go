package kb

import (
	"context"

	"github.com/Abraxas-365/toolkit/pkg/database"
)

type Repository interface {
	GetKnowlegeBaseConfig() (*KnowlegeBaseConfig, error)
	SaveData(ctx context.Context, data DataFile) (*DataFile, error)
	DeleteData(ctx context.Context, dataId int) (*DataFile, error)
	GetData(ctx context.Context, page, pageSize int) (database.PaginatedRecord[DataFile], error)
	GetDataById(ctx context.Context, id int) (*DataFile, error)
}
