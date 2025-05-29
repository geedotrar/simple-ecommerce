package helpers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	client     *s3.Client
	bucketName string
	folder     string
	region     string
}

func NewS3Uploader(bucketName, folder string) (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		client:     client,
		bucketName: bucketName,
		folder:     folder,
		region:     cfg.Region,
	}, nil
}

func (u *S3Uploader) UploadWithFallback(ctx context.Context, fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName := fmt.Sprintf("%d-%s", fileHeader.Size, filepath.Base(fileHeader.Filename))

	localFilePath := filepath.Join(uploadDir, fileName)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create local upload dir: %w", err)
	}

	dst, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return "", fmt.Errorf("failed to save file locally: %w", err)
	}
	dst.Close()

	f, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open local file for upload: %w", err)
	}
	defer f.Close()

	key := fmt.Sprintf("%s/%s", u.folder, fileName)
	uploader := manager.NewUploader(u.client)

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	if err != nil {
		fmt.Printf("S3 upload failed, fallback to local file: %v\n", err)
		return "/uploads/products/" + fileName, nil
	}

	_ = os.Remove(localFilePath)

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucketName, u.region, key)
	return url, nil
}

func (u *S3Uploader) DeleteFileFromS3(ctx context.Context, fileURL string) error {
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", u.bucketName, u.region)
	if !strings.HasPrefix(fileURL, prefix) {
		return fmt.Errorf("file URL does not match S3 bucket format")
	}
	key := strings.TrimPrefix(fileURL, prefix)

	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucketName),
		Key:    aws.String(key),
	})

	return err
}

func IsS3URL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func GetFileNameFromURL(url string) string {
	return filepath.Base(url)
}
