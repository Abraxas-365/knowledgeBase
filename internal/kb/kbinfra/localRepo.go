package kbinfra

import (
	"os"
	"strconv"

	"github.com/Abraxas-365/opd/internal/kb"
	"github.com/Abraxas-365/toolkit/pkg/errors"
)

type LocalConfig struct{}

func NewLocalConfig() *LocalConfig {
	return &LocalConfig{}
}

func (lc *LocalConfig) GetKnowlegeBaseConfig() (*kb.KnowlegeBaseConfig, error) {
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

	return &kb.KnowlegeBaseConfig{
		ID:              id,
		NumberOfResults: numberOfResultsInt,
		Region:          region,
		Model: kb.ModelInformation{
			ModelId: modelId,
			Prompt:  modelPrompt,
		}}, nil

}
