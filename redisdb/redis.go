package redisdb

import (
	"context"
	"crypto/tls"

	"github.com/hyperits/gosuite/logger"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ErrNotConfigured = errors.New("Redis is not configured")

type RedisConfig struct {
	Address           string
	Username          string
	Password          string
	DB                int
	UseTLS            bool
	MasterName        string
	SentinelUsername  string
	SentinelPassword  string
	SentinelAddresses []string
	ClusterAddresses  []string
	MaxRedirects      *int //  for clustererd mode only, number of redirects to follow, defaults to 2
}

type RedisComponent struct {
	conf   *RedisConfig
	client redis.UniversalClient
}

func NewRedisComponent(conf *RedisConfig) (*RedisComponent, error) {
	if !conf.IsConfigured() {
		return nil, ErrNotConfigured
	}

	client, err := newRedisClient(conf)
	if err != nil {
		return nil, err
	}

	return &RedisComponent{
		conf:   conf,
		client: client,
	}, nil
}

func (c *RedisComponent) Client() redis.UniversalClient {
	return c.client
}

func (c *RedisComponent) Config() *RedisConfig {
	return c.conf
}

func (r *RedisConfig) IsConfigured() bool {
	if r.Address != "" {
		return true
	}
	if len(r.SentinelAddresses) > 0 {
		return true
	}
	if len(r.ClusterAddresses) > 0 {
		return true
	}
	return false
}

func (r *RedisConfig) GetMaxRedirects() int {
	if r.MaxRedirects != nil {
		return *r.MaxRedirects
	}
	return 2
}

func newRedisClient(conf *RedisConfig) (redis.UniversalClient, error) {
	if conf == nil {
		return nil, errors.New("redis config is nil")
	}

	if !conf.IsConfigured() {
		return nil, ErrNotConfigured
	}

	var rcOptions *redis.UniversalOptions
	var rc redis.UniversalClient
	var tlsConfig *tls.Config

	if conf.UseTLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if len(conf.SentinelAddresses) > 0 {
		logger.Infof("connecting to redis %v %v addr %v masterName %v", "sentinel", true, conf.SentinelAddresses, conf.MasterName)
		rcOptions = &redis.UniversalOptions{
			Addrs:            conf.SentinelAddresses,
			SentinelUsername: conf.SentinelUsername,
			SentinelPassword: conf.SentinelPassword,
			MasterName:       conf.MasterName,
			Username:         conf.Username,
			Password:         conf.Password,
			DB:               conf.DB,
			TLSConfig:        tlsConfig,
		}
	} else if len(conf.ClusterAddresses) > 0 {
		logger.Infof("connecting to redis %v %v addr %v", "cluster", true, conf.ClusterAddresses)
		rcOptions = &redis.UniversalOptions{
			Addrs:        conf.ClusterAddresses,
			Username:     conf.Username,
			Password:     conf.Password,
			DB:           conf.DB,
			TLSConfig:    tlsConfig,
			MaxRedirects: conf.GetMaxRedirects(),
		}
	} else {
		logger.Infof("connecting to redis %v %v addr %v", "simple", true, conf.Address)
		rcOptions = &redis.UniversalOptions{
			Addrs:     []string{conf.Address},
			Username:  conf.Username,
			Password:  conf.Password,
			DB:        conf.DB,
			TLSConfig: tlsConfig,
		}
	}
	rc = redis.NewUniversalClient(rcOptions)

	if err := rc.Ping(context.Background()).Err(); err != nil {
		err = errors.Wrap(err, "unable to connect to redis")
		return nil, err
	}

	return rc, nil
}
