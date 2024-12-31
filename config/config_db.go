package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/levensspel/go-gin-template/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDbInit() (*gorm.DB, error) {
	dsn := newConfig().GetDsn()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Println("ini error: ", err)
		return nil, err
	}

	err = db.AutoMigrate(
		entity.User{},
	)

	return db, err
}

type dbConfig struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
}

func newConfig() *dbConfig {
	config := dbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DbName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	return &config
}

func (dbConfig *dbConfig) GetDsn() string {

	mode := os.Getenv("MODE")
	var dsn string

	if mode == "production" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.DbName,
			dbConfig.Port,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.DbName,
			dbConfig.Port,
		)
	}

	return dsn
}
