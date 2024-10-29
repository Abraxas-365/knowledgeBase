package kb

type ModelInformation struct {
	ModelId string `json:"modelId"`
	Prompt  string `json:"prompt"`
}

type KnowlegeBaseConfig struct {
	ID              string           `json:"id"`
	NumberOfResults int              `json:"numberOfResults"`
	Region          string           `json:"region"`
	Model           ModelInformation `json:"model"`
}
