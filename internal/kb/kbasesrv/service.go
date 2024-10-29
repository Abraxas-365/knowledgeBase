package kbsrv

import (
	"context"

	"github.com/Abraxas-365/opd/internal/kb"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/Abraxas-365/toolkit/pkg/s3client"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

type Service struct {
	kbClient *bedrockagentruntime.Client
	repo     kb.Repository
	s3Client s3client.Client
}

func New(kbClient *bedrockagentruntime.Client, repo kb.Repository, s3 s3client.Client) *Service {
	return &Service{
		kbClient: kbClient,
		repo:     repo,
		s3Client: s3,
	}
}

func (s *Service) CompleteAnswerWithMetadata(ctx context.Context, userMessage string, sessionID *string) (*bedrockagentruntime.RetrieveAndGenerateOutput, error) {
	kbConf, err := s.repo.GetKnowlegeBaseConfig()
	if err != nil {
		return nil, err
	}

	output, err := s.kbClient.RetrieveAndGenerate(
		context.TODO(),
		&bedrockagentruntime.RetrieveAndGenerateInput{
			SessionId: sessionID,
			Input: &types.RetrieveAndGenerateInput{
				Text: aws.String(userMessage),
			},
			RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
				Type: types.RetrieveAndGenerateTypeKnowledgeBase,
				KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
					KnowledgeBaseId: aws.String(kbConf.ID),
					ModelArn:        aws.String(kbConf.Model.ModelId),
					RetrievalConfiguration: &types.KnowledgeBaseRetrievalConfiguration{
						VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
							NumberOfResults: aws.Int32(int32(kbConf.NumberOfResults)),
						},
					},
					GenerationConfiguration: &types.GenerationConfiguration{
						PromptTemplate: &types.PromptTemplate{
							TextPromptTemplate: aws.String(kbConf.Model.Prompt),
						},
					},
				},
			},
		},
	)
	if err != nil {
		return nil, errors.ErrServiceUnavailable(err.Error())
	}

	return output, nil

}

func (s *Service) GeneratePutURL(key string) (string, error) {
	return s.s3Client.GeneratePresignedPutURL(key, 60)
}

func (s *Service) DeleteObject(key string) error {
	return s.s3Client.DeleteFile(key)
}

func (s *Service) LisObjects(pageSize int32, continuationToken *string) ([]string, *string, error) {
	return s.s3Client.ListFiles(pageSize, continuationToken)
}
