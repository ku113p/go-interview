package minio

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"go-interview/internal/message/domain"
)

var _ domain.FileDeleter = (*MinioRepository)(nil)
var _ domain.FileGetter = (*MinioRepository)(nil)
var _ domain.FileSaver = (*MinioRepository)(nil)

type MinioRepository struct {
	client     *minio.Client
	bucketName string
}

func NewMinioRepository(client *minio.Client, bucketName string) *MinioRepository {
	return &MinioRepository{
		client:     client,
		bucketName: bucketName,
	}
}

func (r *MinioRepository) Delete(ctx context.Context, path string) error {
	return r.client.RemoveObject(ctx, r.bucketName, path, minio.RemoveObjectOptions{})
}

func (r *MinioRepository) Get(ctx context.Context, path string) (io.Reader, error) {
	return r.client.GetObject(ctx, r.bucketName, path, minio.GetObjectOptions{})
}

func (r *MinioRepository) Save(ctx context.Context, stream io.Reader) (*string, error) {
	objectName := uuid.NewString()
	_, err := r.client.PutObject(ctx, r.bucketName, objectName, stream, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &objectName, nil
}
