package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *dbProxy

type dbProxy struct {
	DbConn *gorm.DB
}

func MustInitDb() {
	dbTem, err := initDb()

	if err != nil {
		panic(err)
	}
	Db = &dbProxy{
		DbConn: dbTem,
	}

}

func initDb() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 0,             // 慢 SQL 阈值  todo: 后续改成time.second
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // 禁用彩色打印
		},
	)
	dsn := "root:w123258123@tcp(localhost:3306)/Experiment?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		return nil, err
	}
	return db, nil
}
