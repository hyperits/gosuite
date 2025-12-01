package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/hyperits/gosuite/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 错误定义
var (
	ErrEmptyEndpoint  = errors.New("s3: endpoint is required")
	ErrEmptyAccessKey = errors.New("s3: access key is required")
	ErrEmptySecret    = errors.New("s3: secret is required")
	ErrEmptyBucket    = errors.New("s3: bucket is required")
)

// S3Config S3 客户端配置
type S3Config struct {
	Endpoint       string // S3 服务端点
	AccessKey      string // 访问密钥 ID
	Secret         string // 访问密钥 Secret
	Bucket         string // 默认存储桶名称
	Region         string // 地域
	Secure         bool   // 是否使用 HTTPS
	ForcePathStyle bool   // 是否强制使用路径风格
}

// Validate 验证配置是否有效
func (c *S3Config) Validate() error {
	if c.Endpoint == "" {
		return ErrEmptyEndpoint
	}
	if c.AccessKey == "" {
		return ErrEmptyAccessKey
	}
	if c.Secret == "" {
		return ErrEmptySecret
	}
	if c.Bucket == "" {
		return ErrEmptyBucket
	}
	return nil
}

// S3Client S3 客户端
type S3Client struct {
	client *minio.Client
	config *S3Config
}

// NewS3Client 创建新的 S3 客户端
func NewS3Client(config *S3Config) (*S3Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	var pathStyle minio.BucketLookupType
	if config.ForcePathStyle {
		pathStyle = minio.BucketLookupPath
	} else {
		pathStyle = minio.BucketLookupDNS
	}

	s3Client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(config.AccessKey, config.Secret, ""),
		Secure:       config.Secure,
		Region:       config.Region,
		BucketLookup: pathStyle,
	})
	if err != nil {
		return nil, err
	}

	comp := &S3Client{
		client: s3Client,
		config: config,
	}

	if err := comp.ensureDefaultBucket(context.Background()); err != nil {
		logger.Warnf("failed to ensure default bucket [%s]: %v", config.Bucket, err)
	}

	return comp, nil
}

// ensureDefaultBucket 确保默认 bucket 存在，不存在则创建
func (c *S3Client) ensureDefaultBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.config.Bucket)
	if err != nil {
		return err
	}

	if !exists {
		return c.client.MakeBucket(ctx, c.config.Bucket, minio.MakeBucketOptions{Region: c.config.Region})
	}
	return nil
}

// Client 返回底层 minio 客户端
func (c *S3Client) Client() *minio.Client {
	return c.client
}

// Config 返回当前配置
func (c *S3Client) Config() *S3Config {
	return c.config
}

// resolveBucket 解析 bucket 名称，如果为空则使用默认 bucket
func (c *S3Client) resolveBucket(bucket string) string {
	if bucket == "" {
		return c.config.Bucket
	}
	return bucket
}

// ListObjects 列出指定前缀的对象
func (c *S3Client) ListObjects(ctx context.Context, bucket string, prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	bucket = c.resolveBucket(bucket)

	var objects []minio.ObjectInfo
	objectCh := c.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

// StatObject 获取对象的元信息
func (c *S3Client) StatObject(ctx context.Context, bucket string, objectName string) (minio.ObjectInfo, error) {
	bucket = c.resolveBucket(bucket)
	return c.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
}

// ObjectExists 检查对象是否存在
func (c *S3Client) ObjectExists(ctx context.Context, bucket string, objectName string) (bool, error) {
	_, err := c.StatObject(ctx, bucket, objectName)
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetObject 获取对象并写入目标 writer
func (c *S3Client) GetObject(ctx context.Context, bucket string, objectName string, dst io.Writer) error {
	bucket = c.resolveBucket(bucket)

	object, err := c.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	_, err = io.Copy(dst, object)
	return err
}

// GetObjectAsBytes 获取对象并返回字节数组
func (c *S3Client) GetObjectAsBytes(ctx context.Context, bucket string, objectName string) ([]byte, error) {
	var buf bytes.Buffer
	if err := c.GetObject(ctx, bucket, objectName, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UploadObject 上传对象到指定 bucket
func (c *S3Client) UploadObject(ctx context.Context, bucket string, objectName string, data io.Reader, size int64, contentType string) error {
	bucket = c.resolveBucket(bucket)

	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}

	_, err := c.client.PutObject(ctx, bucket, objectName, data, size, opts)
	return err
}

// UploadObjectFromBytes 从字节数组上传对象
func (c *S3Client) UploadObjectFromBytes(ctx context.Context, bucket string, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	return c.UploadObject(ctx, bucket, objectName, reader, int64(len(data)), contentType)
}

// DeleteObject 删除指定对象
func (c *S3Client) DeleteObject(ctx context.Context, bucket string, objectName string) error {
	bucket = c.resolveBucket(bucket)
	return c.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

// DeleteObjectsByPrefix 删除指定前缀的所有对象
func (c *S3Client) DeleteObjectsByPrefix(ctx context.Context, bucket, prefix string) error {
	bucket = c.resolveBucket(bucket)

	objectCh := c.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var lastErr error
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		if err := c.client.RemoveObject(ctx, bucket, object.Key, minio.RemoveObjectOptions{}); err != nil {
			logger.Errorf("failed to delete object [%s]: %v", object.Key, err)
			lastErr = err
		}
	}
	return lastErr
}

// CopyObject 复制对象到另一个位置
func (c *S3Client) CopyObject(ctx context.Context, srcBucket, srcObject, dstBucket, dstObject string) error {
	srcBucket = c.resolveBucket(srcBucket)
	dstBucket = c.resolveBucket(dstBucket)

	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcObject,
	}
	dst := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: dstObject,
	}

	_, err := c.client.CopyObject(ctx, dst, src)
	return err
}

// PresignedGetURL 生成下载对象的预签名 URL
func (c *S3Client) PresignedGetURL(ctx context.Context, bucket string, objectName string, expires time.Duration) (*url.URL, error) {
	bucket = c.resolveBucket(bucket)
	return c.client.PresignedGetObject(ctx, bucket, objectName, expires, nil)
}

// PresignedPutURL 生成上传对象的预签名 URL
func (c *S3Client) PresignedPutURL(ctx context.Context, bucket string, objectName string, expires time.Duration) (*url.URL, error) {
	bucket = c.resolveBucket(bucket)
	return c.client.PresignedPutObject(ctx, bucket, objectName, expires)
}

// BucketExists 检查 bucket 是否存在
func (c *S3Client) BucketExists(ctx context.Context, bucket string) (bool, error) {
	bucket = c.resolveBucket(bucket)
	return c.client.BucketExists(ctx, bucket)
}

// MakeBucket 创建新的 bucket
func (c *S3Client) MakeBucket(ctx context.Context, bucket string) error {
	return c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: c.config.Region})
}

// RemoveBucket 删除空的 bucket
func (c *S3Client) RemoveBucket(ctx context.Context, bucket string) error {
	return c.client.RemoveBucket(ctx, bucket)
}
