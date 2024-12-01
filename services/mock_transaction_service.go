// services/mock_transaction_service.go
package services

import (
	"errors"
	"sync"

	"crypto-exchange/models"
)

// TransactionService defines the methods for transaction operations.
type TransactionService interface {
	CreateTransaction(tx models.Transaction) (models.Transaction, error)
	GetTransactionByID(id string) (models.Transaction, error)
}

// MockTransactionService is a mock implementation of TransactionService.
type MockTransactionService struct {
	transactions map[string]models.Transaction
	mutex        sync.RWMutex
}

// NewMockTransactionService creates a new instance of MockTransactionService.
func NewMockTransactionService() TransactionService {
	return &MockTransactionService{
		transactions: make(map[string]models.Transaction),
	}
}

// CreateTransaction adds a new transaction to the mock store.
func (s *MockTransactionService) CreateTransaction(tx models.Transaction) (models.Transaction, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if tx.ID == "" {
		return models.Transaction{}, errors.New("transaction ID cannot be empty")
	}
	if _, exists := s.transactions[tx.ID]; exists {
		return models.Transaction{}, errors.New("transaction ID already exists")
	}
	s.transactions[tx.ID] = tx
	return tx, nil
}

// GetTransactionByID retrieves a transaction by ID from the mock store.
func (s *MockTransactionService) GetTransactionByID(id string) (models.Transaction, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	tx, exists := s.transactions[id]
	if !exists {
		return models.Transaction{}, errors.New("transaction not found")
	}
	return tx, nil
}