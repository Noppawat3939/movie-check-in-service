package postgresl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SSLMode  string
}

func NewDB() (*gorm.DB, error) {

	dsn := buildDSN(Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_EXTERNAL_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.New("failed to open db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.New("failed to get sql db")
	}

	// pool setting
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	// verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres %w", err)
	}

	log.Print("connected to database")

	return db, nil
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("failed to get sql db for closing: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("failed to close db: %v", err)
	}
}

func buildDSN(cfg Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode)
}
