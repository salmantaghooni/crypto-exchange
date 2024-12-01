// services/transaction_service.go
package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"crypto-exchange/models"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TransactionServiceDB combines all services for transaction operations.
type TransactionServiceDB struct {
	DB           *gorm.DB
	Logger       zerolog.Logger
	RedisService *RedisService
	CassandraSvc *CassandraService
	KafkaService *KafkaService
}

// NewTransactionService initializes a new TransactionServiceDB.
func NewTransactionService(db *gorm.DB, logger zerolog.Logger, redisSvc *RedisService, cassandraSvc *CassandraService, kafkaSvc *KafkaService) *TransactionServiceDB {
	return &TransactionServiceDB{
		DB:           db,
		Logger:       logger,
		RedisService: redisSvc,
		CassandraSvc: cassandraSvc,
		KafkaService: kafkaSvc,
	}
}

// CreateTransaction handles the creation of a new transaction across multiple services.
func (s *TransactionServiceDB) CreateTransaction(tx models.Transaction) (models.Transaction, error) {
	// Start a database transaction
	txDB := s.DB.Begin()
	if txDB.Error != nil {
		s.Logger.Error().Err(txDB.Error).Msg("Failed to start DB transaction")
		return models.Transaction{}, txDB.Error
	}

	// Create the transaction in PostgreSQL
	if err := txDB.Create(&tx).Error; err != nil {
		s.Logger.Error().Err(err).Msg("Failed to create transaction in PostgreSQL")
		txDB.Rollback()
		return models.Transaction{}, err
	}

	// Insert into Cassandra
	if err := s.CassandraSvc.InsertTransaction(tx); err != nil {
		s.Logger.Error().Err(err).Msg("Failed to create transaction in Cassandra")
		txDB.Rollback()
		return models.Transaction{}, err
	}

	// Cache the transaction in Redis
	ctx := context.Background()
	txJSON, err := json.Marshal(tx)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to marshal transaction for Redis")
		txDB.Rollback()
		return models.Transaction{}, err
	}

	if err := s.RedisService.Set(ctx, tx.ID, string(txJSON), time.Minute*10); err != nil {
		s.Logger.Error().Err(err).Msg("Failed to set transaction in Redis")
		// Not rolling back as caching is not critical
	}

	// Publish the transaction event to Kafka
	if err := s.KafkaService.Publish(string(txJSON)); err != nil {
		s.Logger.Error().Err(err).Msg("Failed to publish transaction to Kafka")
		// Not rolling back as event publishing is not critical
	}

	// Commit the transaction
	if err := txDB.Commit().Error; err != nil {
		s.Logger.Error().Err(err).Msg("Failed to commit DB transaction")
		return models.Transaction{}, err
	}

	s.Logger.Info().
		Str("transaction_id", tx.ID).
		Msg("Transaction created successfully across all services")

	return tx, nil
}

// GetTransactionByID retrieves a transaction by ID, utilizing Redis cache.
func (s *TransactionServiceDB) GetTransactionByID(id string) (models.Transaction, error) {
	// Attempt to retrieve from Redis cache
	ctx := context.Background()
	cachedTx, err := s.RedisService.Get(ctx, id)
	if err == nil {
		var tx models.Transaction
		if err := json.Unmarshal([]byte(cachedTx), &tx); err == nil {
			s.Logger.Info().
				Str("transaction_id", id).
				Msg("Transaction retrieved from Redis cache")
			return tx, nil
		}
	}

	// If not in cache, retrieve from PostgreSQL
	var tx models.Transaction
	if err := s.DB.First(&tx, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warn().
				Str("transaction_id", id).
				Msg("Transaction not found in PostgreSQL")
			return models.Transaction{}, errors.New("transaction not found")
		}
		s.Logger.Error().Err(err).Msg("Failed to retrieve transaction from PostgreSQL")
		return models.Transaction{}, err
	}

	// Cache the retrieved transaction in Redis
	txJSON, err := json.Marshal(tx)
	if err == nil {
		s.RedisService.Set(ctx, id, string(txJSON), time.Minute*10)
	}

	s.Logger.Info().
		Str("transaction_id", id).
		Msg("Transaction retrieved from PostgreSQL")

	return tx, nil
}