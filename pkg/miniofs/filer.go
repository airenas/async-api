package miniofs

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/google/uuid"
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

// NewFiler creates Minio file saver
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
	goapp.Log.Info().Str("file", name).Int64("size b", info.Size).Msgf("saved")
	return nil
}

// LoadFile loads file from s3/minio
func (fs *Filer) LoadFile(ctx context.Context, name string) (io.ReadSeekCloser, error) {
	goapp.Log.Info().Str("file", name).Msg("load")
	res, err := fs.minioClient.GetObject(ctx, fs.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("can't save %s: %w", name, err)
	}
	return &fileWrap{f: res}, nil
}

// Clean removes all files from s3/minio starting by prefix
func (fs *Filer) Clean(ctx context.Context, prefix string) error {
	if prefix == "" {
		return fmt.Errorf("no prefix")
	}
	_, err := uuid.Parse(prefix)
	if err != nil {
		return fmt.Errorf("wrong ID")
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	goapp.Log.Info().Str("prefix", prefix).Msg("clean fs")
	objectCh := fs.minioClient.ListObjects(ctx, fs.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	rmChan := fs.minioClient.RemoveObjectsWithResult(ctx, fs.bucket, objectCh, minio.RemoveObjectsOptions{GovernanceBypass: true})

	for res := range rmChan {
		if res.Err != nil {
			return fmt.Errorf("can't remove %s: %w", res.ObjectName, res.Err)
		}
		goapp.Log.Info().Str("file", res.ObjectName).Msg("removed")
	}
	return nil
}

type fileWrap struct {
	f *minio.Object
}

// Read implements io.ReadSeekCloser
func (fw *fileWrap) Read(p []byte) (n int, err error) {
	return fw.f.Read(p)
}

// Seek implements io.ReadSeekCloser
func (fw *fileWrap) Seek(offset int64, whence int) (int64, error) {
	return fw.f.Seek(offset, whence)
}

// Close implements io.ReadSeekCloser
func (fw *fileWrap) Close() error {
	return fw.f.Close()
}

// Stat returns file stat
func (fw *fileWrap) Stat() (fs.FileInfo, error) {
	st, err := fw.f.Stat()
	if err != nil {
		return nil, err
	}
	return &statsWrap{oi: st}, nil
}

type statsWrap struct {
	oi minio.ObjectInfo
}

// IsDir implements fs.FileInfo
func (sw *statsWrap) IsDir() bool {
	return false
}

// ModTime implements fs.FileInfo
func (sw *statsWrap) ModTime() time.Time {
	return sw.oi.LastModified
}

// Mode implements fs.FileInfo
func (sw *statsWrap) Mode() fs.FileMode {
	return fs.ModeTemporary
}

// Name implements fs.FileInfo
func (sw *statsWrap) Name() string {
	return sw.oi.Key
}

// Size implements fs.FileInfo
func (sw *statsWrap) Size() int64 {
	return sw.oi.Size
}

// Sys implements fs.FileInfo
func (sw *statsWrap) Sys() any {
	return nil
}
