package postgresdb

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

type PostgresComponent struct {
	db   *gorm.DB
	conf *PostgresConfig
}

func NewPostgresComponent(conf *PostgresConfig) (*PostgresComponent, error) {
	db, err := newPostgres(conf)
	if err != nil {
		return nil, err
	}

	return &PostgresComponent{
		db:   db,
		conf: conf,
	}, nil
}

func (c *PostgresComponent) DB() *gorm.DB {
	return c.db
}

func (c *PostgresComponent) Config() *PostgresConfig {
	return c.conf
}

func newPostgres(conf *PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
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
