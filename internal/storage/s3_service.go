package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type s3StorageService struct {
	client     *s3.Client
	bucketName string
}

// NewS3StorageService constructs an s3StorageService.
// Bucket name is read from AWS_S3_BUCKET env var; defaults to "myawsbucketelito".
func NewS3StorageService(client *s3.Client) StorageService {
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		bucket = "myawsbucketelito"
	}
	return &s3StorageService{client: client, bucketName: bucket}
}

func (s *s3StorageService) Upload(ctx context.Context, name string, content []byte) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(name),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return "", fmt.Errorf("%w: uploading %q: %s", ErrStorage, name, err)
	}
	return "file uploaded: " + name, nil
}

func (s *s3StorageService) Download(ctx context.Context, name string) ([]byte, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(name),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: downloading %q: %s", ErrStorage, name, err)
	}
	defer output.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(output.Body); err != nil {
		return nil, fmt.Errorf("%w: reading body for %q: %s", ErrStorage, name, err)
	}
	return buf.Bytes(), nil
}

func (s *s3StorageService) Delete(ctx context.Context, name string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("%w: deleting %q: %s", ErrStorage, name, err)
	}
	return nil
}

func (s *s3StorageService) List(ctx context.Context) ([]string, error) {
	output, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: listing objects: %s", ErrStorage, err)
	}

	names := make([]string, len(output.Contents))
	for i, obj := range output.Contents {
		names[i] = aws.ToString(obj.Key)
	}
	return names, nil
}

func isNotFound(err error) bool {
	var noSuchKey *types.NoSuchKey
	if smithyErr, ok := err.(*smithy.GenericAPIError); ok {
		return smithyErr.Code == "NoSuchKey"
	}
	return errors.As(err, &noSuchKey)
}
