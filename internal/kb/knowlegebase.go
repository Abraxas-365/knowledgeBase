package kb

type ModelInformation struct {
	ModelId string `json:"modelId"`
	Prompt  string `json:"prompt"`
}

type KnowlegeBaseConfig struct {
	ID              string           `json:"id"`
	S3DataSurce     string           `json:"s3DataSurce"`
	NumberOfResults int              `json:"numberOfResults"`
	Region          string           `json:"region"`
	Model           ModelInformation `json:"model"`
}

type DataFile struct {
	ID        int    `db:"id" json:"id"`
	Filename  string `db:"filename" json:"filename"`
	S3Key     string `db:"s3_key" json:"s3_key"`
	UserID    string `db:"user_id" json:"user_id"`
	UserEmail string `db:"user_email" json:"user_email"`
}
