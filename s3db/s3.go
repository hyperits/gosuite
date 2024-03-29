package s3db

import (
	"context"
	"io"
	"strings"

	"github.com/hyperits/gosuite/logger"
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

type S3Component struct {
	client *minio.Client
	config *S3Config
}

func NewS3Component(config *S3Config) (*S3Component, error) {
	// Initialize minio client object.
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

	comp := &S3Component{
		client: s3Client,
		config: config,
	}

	comp.makeDefaultBucket()

	return comp, nil
}

func (c *S3Component) makeDefaultBucket() {
	exists, err := c.client.BucketExists(context.Background(), c.config.Bucket)
	if err != nil {
		logger.Errorf("error check bucket exists [%s], [%s]", c.config.Bucket, err.Error())
		return
	}

	if !exists {
		err := c.client.MakeBucket(context.Background(), c.config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Errorf("error create bucket [%s], [%s]", c.config.Bucket, err.Error())
		}
	}
}

func (c *S3Component) Client() *minio.Client {
	return c.client
}

func (c *S3Component) Config() *S3Config {
	return c.config
}

func (c *S3Component) ListObjects(bucket string, prefix string, recursive bool) ([]minio.ObjectInfo, error) {
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
		if !strings.HasSuffix(object.Key, ".mp4") {
			// ignore json
			continue
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func (c *S3Component) GetObject(bucket string, objectName string, dst io.Writer) error {
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
