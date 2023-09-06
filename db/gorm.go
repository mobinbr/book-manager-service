package db

import (
	"BookManager/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormDB represents a connection to the database
type GormDB struct {
	cfg config.Config
	db  *gorm.DB
}

// NewGormDB creates a new GormDB instance using the provided configuration
func NewGormDB(cfg config.Config) (*GormDB, error) {
	c := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	// Create a new connection
	db, err := gorm.Open(postgres.Open(c), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormDB{
		cfg: cfg,
		db:  db,
	}, nil
}

// CreateSchemas creates database schema for models
func (gdb *GormDB) CreateSchemas() error {
	// Auto migrate User Accounts & Books 
	err := gdb.db.AutoMigrate(&User{}, &Book{})
	if err != nil {
		return err
	}

	return nil
}
