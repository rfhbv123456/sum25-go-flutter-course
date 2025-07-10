package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds database configuration
type Config struct {
	DatabasePath    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns a default database configuration
func DefaultConfig() *Config {
	return &Config{
		DatabasePath:    "./lab04.db",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}
}

// InitDB инициализирует соединение с базой данных с настройками по умолчанию
func InitDB() (*sql.DB, error) {
	return InitDBWithConfig(DefaultConfig())
}

// InitDBWithConfig инициализирует соединение с базой данных с кастомными настройками
func InitDBWithConfig(config *Config) (*sql.DB, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	// Открываем соединение
	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return nil, err
	}
	// Настройки пула соединений
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	// Проверяем соединение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// CloseDB закрывает соединение с базой данных
func CloseDB(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	return db.Close()
}
