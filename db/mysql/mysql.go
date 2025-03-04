package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MySQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

// MySQLClient 提供MySQL数据库连接和操作的客户端
type MySQLClient struct {
	db   *gorm.DB
	conf *MySQLConfig
}

// NewMySQLClient 创建一个新的MySQL客户端实例
func NewMySQLClient(conf *MySQLConfig) (*MySQLClient, error) {
	db, err := connectMySQL(conf)
	if err != nil {
		return nil, err
	}

	return &MySQLClient{
		db:   db,
		conf: conf,
	}, nil
}

// DB 返回底层的gorm数据库连接
func (c *MySQLClient) DB() *gorm.DB {
	return c.db
}

// Config 返回MySQL配置信息
func (c *MySQLClient) Config() *MySQLConfig {
	return c.conf
}

// connectMySQL 连接到MySQL数据库并返回数据库连接
func connectMySQL(conf *MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect mysql failed: %w", err)
	}

	return db, nil
}
