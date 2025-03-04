package postgres

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	SSLMode  string
	TimeZone string
}

type PostgresClient struct {
	db   *gorm.DB
	conf *PostgresConfig
}

func NewPostgresClient(conf *PostgresConfig) (*PostgresClient, error) {
	db, err := newPostgres(conf)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{
		db:   db,
		conf: conf,
	}, nil
}

func (c *PostgresClient) DB() *gorm.DB {
	return c.db
}

func (c *PostgresClient) Config() *PostgresConfig {
	return c.conf
}

func newPostgres(conf *PostgresConfig) (*gorm.DB, error) {
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
		log.Fatal(err.Error())
	}

	return db, err
}
