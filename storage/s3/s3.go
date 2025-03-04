package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/hyperits/gosuite/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	Endpoint       string
	AccessKey      string
	Secret         string
	Bucket         string
	Region         string
	Secure         bool
	ForcePathStyle bool
}

type S3Client struct {
	client *minio.Client
	config *S3Config
}

func NewS3Client(config *S3Config) (*S3Client, error) {
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

	comp.makeDefaultBucket()

	return comp, nil
}

func (c *S3Client) makeDefaultBucket() {
	exists, err := c.client.BucketExists(context.Background(), c.config.Bucket)
	if err != nil {
		log.Errorf("error checking bucket exists [%s], [%s]", c.config.Bucket, err.Error())
		return
	}

	if !exists {
		err := c.client.MakeBucket(context.Background(), c.config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Errorf("error creating bucket [%s], [%s]", c.config.Bucket, err.Error())
		}
	}
}

func (c *S3Client) Client() *minio.Client {
	return c.client
}

func (c *S3Client) Config() *S3Config {
	return c.config
}

// ListObjects lists objects with a specified prefix and recursive flag
func (c *S3Client) ListObjects(bucket string, prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if bucket == "" {
		bucket = c.config.Bucket
	}

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

// GetObject fetches an object and writes it to a destination writer
func (c *S3Client) GetObject(bucket string, objectName string, dst io.Writer) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	object, err := c.client.GetObject(context.Background(), bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	if _, err = io.Copy(dst, object); err != nil {
		return err
	}
	return nil
}

// UploadObject uploads an object to the specified bucket
func (c *S3Client) UploadObject(bucket string, objectName string, data io.Reader, size int64, contentType string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	_, err := c.client.PutObject(context.Background(), bucket, objectName, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// DeleteObject deletes an object from the specified bucket
func (c *S3Client) DeleteObject(bucket string, objectName string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	return c.client.RemoveObject(context.Background(), bucket, objectName, minio.RemoveObjectOptions{})
}

// CopyObject copies an object to another location
func (c *S3Client) CopyObject(srcBucket, srcObject, dstBucket, dstObject string) error {
	if srcBucket == "" {
		srcBucket = c.config.Bucket
	}
	if dstBucket == "" {
		dstBucket = c.config.Bucket
	}

	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcObject,
	}
	dst := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: dstObject,
	}
	_, err := c.client.CopyObject(context.Background(), dst, src)
	return err
}

// DeleteObjectsByPrefix deletes all objects with a specific prefix
func (c *S3Client) DeleteObjectsByPrefix(bucket, prefix string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	ctx := context.Background()

	objectCh := c.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		err := c.client.RemoveObject(ctx, bucket, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			log.Errorf("error deleting object [%s], [%s]", object.Key, err.Error())
		}
	}
	return nil
}

// UploadObjectFromBytes uploads an object using a byte slice as source
func (c *S3Client) UploadObjectFromBytes(bucket string, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	return c.UploadObject(bucket, objectName, reader, int64(len(data)), contentType)
}
