// models/transaction.go
package models

type Transaction struct {
    ID     uint  `json:"id" gorm:"primaryKey" binding:"required"`
    UserID uint `json:"user_id" binding:"required"`
    Amount float64 `json:"amount" binding:"required,gt=0"`
    Type   string  `json:"type" binding:"required,oneof=deposit withdrawal"`
    Status string  `json:"status" binding:"required,oneof=pending completed failed"`
    CryptoType string `json:"crypto_type" binding:"required"`
    TransactionID  string   `json:"transaction_id" binding:"required"`
    CryptoAmount   float64 `json:"crypto_amount" binding:"required,gt=0"`
    CryptoSymbol   string  `json:"crypto_symbol" binding:"required"` // e.g., BTC, ETH
    TransactionFee float64 `json:"transaction_fee" binding:"required,gt=0"`
    CreatedAt   int64   `json:"created_at"`
    UpdatedAt   int64   `json:"updated_at"`
    DeletedAt   int64   `json:"deleted_at,omitempty"`
}