package s3

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	AccountID     string
	AccessKey     string
	SecretKey     string
	BucketName    string
	Uploader      *manager.Uploader
	PresignClient *s3.PresignClient
}

type Uploader interface {
	Upload(src io.Reader, bucket string, filename string) error
}

func NewS3Uploader(accountID, accessKey, secretKey string, bucket string) (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	uploader := manager.NewUploader(client)

	presignClient := s3.NewPresignClient(client)

	return &S3Uploader{
		AccountID:     accountID,
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		BucketName:    bucket,
		Uploader:      uploader,
		PresignClient: presignClient,
	}, nil
}

func (s *S3Uploader) Upload(f multipart.File, key string, contentType string) error {
	_, err := s.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &s.BucketName,
		Key:         aws.String(key),
		Body:        f,
		ContentType: &contentType,
	})
	if err != nil {
		return err
	}
	return nil
}
