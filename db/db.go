// Package db 提供数据库客户端的公共接口定义
package db

import (
	"context"
	"io"
)

// Client 数据库客户端通用接口
// 所有数据库客户端（MySQL、PostgreSQL、Redis）都应实现此接口
type Client interface {
	io.Closer

	// Ping 测试数据库连接
	Ping(ctx context.Context) error

	// IsConnected 检查是否已连接
	IsConnected() bool
}

// SQLClient SQL 数据库客户端接口
// MySQL 和 PostgreSQL 客户端实现此接口
type SQLClient interface {
	Client

	// Exec 执行 SQL 语句
	Exec(ctx context.Context, sql string, args ...interface{}) error

	// QueryRow 查询单行
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row

	// Query 查询多行
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)

	// Begin 开始事务
	Begin(ctx context.Context) (Tx, error)
}

// Row 单行查询结果接口
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows 多行查询结果接口
type Rows interface {
	io.Closer
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

// Tx 事务接口
type Tx interface {
	Commit() error
	Rollback() error
	Exec(ctx context.Context, sql string, args ...interface{}) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
}

// KVClient 键值存储客户端接口
// Redis 等键值存储实现此接口
type KVClient interface {
	Client

	// Get 获取值
	Get(ctx context.Context, key string) (string, error)

	// Set 设置值
	Set(ctx context.Context, key string, value interface{}) error

	// SetWithTTL 设置值并指定过期时间（秒）
	SetWithTTL(ctx context.Context, key string, value interface{}, ttlSeconds int) error

	// Del 删除键
	Del(ctx context.Context, keys ...string) error

	// Exists 检查键是否存在
	Exists(ctx context.Context, keys ...string) (int64, error)
}
