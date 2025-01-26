package awssqs

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type SqsManager struct {
	service *sqs.Client
	groupId string
}

func New() (*SqsManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to load configuration: %s", err))
	}
	svc := sqs.NewFromConfig(cfg)

	return &SqsManager{svc, uuid.New().String()}, nil
}

func (mngr *SqsManager) SendMessage(queue, message string) error {
	_, err := mngr.service.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:       aws.String(queue),
		MessageBody:    aws.String(message),
		MessageGroupId: aws.String(mngr.groupId),
	})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send message to broker: %s", err))
	}

	return nil
}
