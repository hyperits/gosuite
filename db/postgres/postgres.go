package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hyperits/gosuite/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Config PostgreSQL 数据库配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	SSLMode  string
	TimeZone string

	// 连接池配置
	MaxOpenConns    int           // 最大打开连接数，默认 25
	MaxIdleConns    int           // 最大空闲连接数，默认 10
	ConnMaxLifetime time.Duration // 连接最大生命周期，默认 5 分钟
	ConnMaxIdleTime time.Duration // 空闲连接最大生命周期，默认 5 分钟
}

// Client PostgreSQL 数据库客户端
type Client struct {
	db     *gorm.DB
	conf   *Config
	closed bool
	mu     sync.RWMutex
}

// NewClient 创建一个新的 PostgreSQL 客户端实例
func NewClient(conf *Config) (*Client, error) {
	if conf == nil {
		return nil, errors.ErrNilConfig
	}

	db, err := connect(conf)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:   db,
		conf: conf,
	}, nil
}

// DB 返回底层的 gorm 数据库连接
func (c *Client) DB() *gorm.DB {
	return c.db
}

// GetConfig 返回 PostgreSQL 配置信息
func (c *Client) GetConfig() *Config {
	return c.conf
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.ErrAlreadyClosed
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return errors.Wrap(err, "postgres.close")
	}

	c.closed = true
	return sqlDB.Close()
}

// Ping 测试数据库连接
func (c *Client) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return errors.ErrAlreadyClosed
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return errors.Wrap(err, "postgres.ping")
	}

	return sqlDB.PingContext(ctx)
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}

// connect 连接到 PostgreSQL 数据库并返回数据库连接
func connect(conf *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		conf.Host,
		conf.Port,
		conf.Username,
		conf.Password,
		conf.DbName,
		conf.SSLMode,
		conf.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect postgres failed: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB failed: %w", err)
	}

	// 设置连接池参数（使用默认值或配置值）
	maxOpenConns := conf.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 25
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)

	maxIdleConns := conf.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)

	connMaxLifetime := conf.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 5 * time.Minute
	}
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	connMaxIdleTime := conf.ConnMaxIdleTime
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = 5 * time.Minute
	}
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	return db, nil
}
