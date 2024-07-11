package mysqldb

import (
	"fmt"
	"log"

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

type MySQLComponent struct {
	db   *gorm.DB
	conf *MySQLConfig
}

func NewMySQLComponent(conf *MySQLConfig) (*MySQLComponent, error) {
	db, err := newMySQL(conf)
	if err != nil {
		return nil, err
	}

	return &MySQLComponent{
		db:   db,
		conf: conf,
	}, nil
}

func (c *MySQLComponent) DB() *gorm.DB {
	return c.db
}

func (c *MySQLComponent) Config() *MySQLConfig {
	return c.conf
}

func newMySQL(conf *MySQLConfig) (*gorm.DB, error) {
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
		log.Fatal(err.Error())
	}

	return db, err
}
