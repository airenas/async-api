package miniofs

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Filer saves files on s3/minio
type Filer struct {
	minioClient *minio.Client
	bucket      string
}

// Options is minio client initializatoin options
type Options struct {
	URL, User, Key, Bucket string
}

//NewFiler creates Minio file saver
func NewFiler(ctx context.Context, opt Options) (*Filer, error) {
	goapp.Log.Info().Msgf("Init MinIO File Storage at: %s(%s)", opt.URL, opt.Bucket)
	if err := validate(opt); err != nil {
		return nil, err
	}
	minioClient, err := minio.New(opt.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(opt.User, opt.Key, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("can't init minio client: %w", err)
	}

	err = minioClient.MakeBucket(ctx, opt.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, opt.Bucket)
		if !(errBucketExists == nil && exists) {
			return nil, fmt.Errorf("can't init bucket: %w", err)
		}
	}
	return &Filer{minioClient: minioClient, bucket: opt.Bucket}, nil
}

func validate(opt Options) error {
	if opt.URL == "" {
		return fmt.Errorf("no URL")
	}
	if opt.User == "" {
		return fmt.Errorf("no user")
	}
	if opt.Bucket == "" {
		return fmt.Errorf("no bucket")
	}
	return nil
}

// SaveFile saves file to s3/minio
func (fs *Filer) SaveFile(ctx context.Context, name string, reader io.Reader) error {
	if strings.Contains(name, "..") {
		return fmt.Errorf("wrong path '%s'", name)
	}
	info, err := fs.minioClient.PutObject(ctx, fs.bucket, name, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("can't save %s: %w", name, err)
	}
	goapp.Log.Info().Msgf("Saved file %s. Size = %s b", name, strconv.FormatInt(info.Size, 10))
	return nil
}
