package awss3

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Manager struct {
	service *s3.Client
}

func New() (*s3Manager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to load configuration: %s", err))
	}
	svc := s3.NewFromConfig(cfg)

	return &s3Manager{svc}, nil
}

func (mngr *s3Manager) DownloadFile(repository, fileName string) (io.ReadCloser, error) {
	result, err := mngr.service.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(repository),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get file %s from s3: %s", fileName, err))
	}

	return result.Body, nil
}

func (mngr *s3Manager) UploadFile(repository, fileName string, data io.Reader) error {
	uploader := manager.NewUploader(mngr.service)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(repository),
		Key:    aws.String(fileName),
		Body:   data,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload file %s to s3: %s", fileName, err))
	}

	return nil
}

func (mngr *s3Manager) DeleteFile(repository, fileName string) error {
	_, err := mngr.service.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(repository),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete file %s from s3: %s", fileName, err))
	}

	return nil
}
