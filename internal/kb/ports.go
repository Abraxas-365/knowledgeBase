package kb

type Repository interface {
	GetKnowlegeBaseConfig() (*KnowlegeBaseConfig, error)
}
