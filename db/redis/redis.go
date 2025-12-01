package redis

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/hyperits/gosuite/errors"
	"github.com/hyperits/gosuite/log"
	"github.com/redis/go-redis/v9"
)

// Config Redis 配置
type Config struct {
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
	MaxRedirects      *int // 仅集群模式，重定向次数，默认 2
}

// Client Redis 客户端
type Client struct {
	conf   *Config
	client redis.UniversalClient
	closed bool
	mu     sync.RWMutex
}

// NewClient 创建 Redis 客户端
func NewClient(conf *Config) (*Client, error) {
	if conf == nil {
		return nil, errors.ErrNilConfig
	}

	if !conf.IsConfigured() {
		return nil, errors.ErrNotConfigured
	}

	client, err := connect(conf)
	if err != nil {
		return nil, err
	}

	return &Client{
		conf:   conf,
		client: client,
	}, nil
}

// UniversalClient 返回底层的 Redis 客户端
func (c *Client) UniversalClient() redis.UniversalClient {
	return c.client
}

// GetConfig 返回 Redis 配置
func (c *Client) GetConfig() *Config {
	return c.conf
}

// Close 关闭 Redis 连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.ErrAlreadyClosed
	}

	c.closed = true
	return c.client.Close()
}

// Ping 测试 Redis 连接
func (c *Client) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return errors.ErrAlreadyClosed
	}

	return c.client.Ping(ctx).Err()
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.client.Ping(ctx).Err() == nil
}

// Get 获取值
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置值（无过期时间）
func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, key, value, 0).Err()
}

// SetWithTTL 设置值并指定过期时间
func (c *Client) SetWithTTL(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	return c.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

// Del 删除键
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// IsConfigured 检查配置是否有效
func (r *Config) IsConfigured() bool {
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

// GetMaxRedirects 获取最大重定向次数
func (r *Config) GetMaxRedirects() int {
	if r.MaxRedirects != nil {
		return *r.MaxRedirects
	}
	return 2
}

// connect 连接到 Redis
func connect(conf *Config) (redis.UniversalClient, error) {
	var rcOptions *redis.UniversalOptions
	var tlsConfig *tls.Config

	if conf.UseTLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	switch {
	case len(conf.SentinelAddresses) > 0:
		log.Infof("connecting to redis sentinel, addr: %v, master: %v", conf.SentinelAddresses, conf.MasterName)
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
	case len(conf.ClusterAddresses) > 0:
		log.Infof("connecting to redis cluster, addr: %v", conf.ClusterAddresses)
		rcOptions = &redis.UniversalOptions{
			Addrs:        conf.ClusterAddresses,
			Username:     conf.Username,
			Password:     conf.Password,
			DB:           conf.DB,
			TLSConfig:    tlsConfig,
			MaxRedirects: conf.GetMaxRedirects(),
		}
	default:
		log.Infof("connecting to redis standalone, addr: %v", conf.Address)
		rcOptions = &redis.UniversalOptions{
			Addrs:     []string{conf.Address},
			Username:  conf.Username,
			Password:  conf.Password,
			DB:        conf.DB,
			TLSConfig: tlsConfig,
		}
	}

	rc := redis.NewUniversalClient(rcOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "redis.connect")
	}

	return rc, nil
}
