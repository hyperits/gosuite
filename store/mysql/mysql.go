package model

import (
	"fmt"
	"log"

	"github.com/hyperits/gosuite/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(conf *config.MySQLConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, err
}
