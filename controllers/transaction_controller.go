// controllers/transaction_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"crypto-exchange/models"
	"crypto-exchange/services"
)

// TransactionController handles transaction-related HTTP requests.
type TransactionController struct {
	Service services.TransactionService
	Logger  zerolog.Logger
}

// NewTransactionController creates a new instance of TransactionController.
func NewTransactionController(service services.TransactionService, logger zerolog.Logger) *TransactionController {
	return &TransactionController{
		Service: service,
		Logger:  logger,
	}
}

// CreateTransaction handles the creation of a new transaction.
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var tx models.Transaction
	// Bind JSON input to Transaction model
	if err := c.ShouldBindJSON(&tx); err != nil {
		tc.Logger.Error().
			Err(err).
			Msg("Invalid transaction payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Create the transaction using the service
	createdTx, err := tc.Service.CreateTransaction(tx)
	if err != nil {
		tc.Logger.Error().
			Err(err).
			Str("transaction_id", tx.ID).
			Msg("Failed to create transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Respond with the created transaction
	tc.Logger.Info().
		Str("transaction_id", createdTx.ID).
		Msg("Transaction created successfully")
	c.JSON(http.StatusCreated, createdTx)
}

// GetTransaction handles fetching a transaction by ID.
func (tc *TransactionController) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		tc.Logger.Warn().Msg("Transaction ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	// Retrieve the transaction using the service
	tx, err := tc.Service.GetTransactionByID(id)
	if err != nil {
		tc.Logger.Error().
			Err(err).
			Str("transaction_id", id).
			Msg("Transaction not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Respond with the retrieved transaction
	tc.Logger.Info().
		Str("transaction_id", tx.ID).
		Msg("Transaction retrieved successfully")
	c.JSON(http.StatusOK, tx)
}