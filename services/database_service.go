// services/database_service.go
package services

import (
	"fmt"

	"crypto-exchange/models"
	"crypto-exchange/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseService encapsulates the GORM DB instance.
type DatabaseService struct {
	DB *gorm.DB
}

// NewDatabaseService initializes the DatabaseService with PostgreSQL.
func NewDatabaseService(cfg config.DatabaseConfig) (*DatabaseService, error) {
	var dsn string
	switch cfg.Type {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.DBName,
			cfg.Postgres.SSLMode,
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		// Automigrate the Transaction model to ensure the table exists.
		if err := db.AutoMigrate(&models.Transaction{}); err != nil {
			return nil, err
		}
		return &DatabaseService{DB: db}, nil
	case "mysql":
		// MySQL implementation can be added here.
		return nil, fmt.Errorf("MySQL not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}