// services/cassandra_service.go
package services

import (
	"crypto-exchange/config"
	"fmt"

	"github.com/gocql/gocql"
)

// CassandraService encapsulates the Cassandra session.
type CassandraService struct {
	Session *gocql.Session
}

// NewCassandraService initializes the CassandraService.
func NewCassandraService(cfg config.CassandraConfig) (*CassandraService, error) {
	cluster := gocql.NewCluster(cfg.Host)
	cluster.Port = cfg.Port
	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %v", err)
	}

	// Create the "transactions" table if it doesn't exist.
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS transactions (
			id text PRIMARY KEY,
			amount double,
			type text,
			status text
		)
	`).Exec()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create transactions table: %v", err)
	}

	return &CassandraService{
		Session: session,
	}, nil
}

// Close terminates the Cassandra session.
func (c *CassandraService) Close() {
	c.Session.Close()
}

// InsertTransaction inserts a new transaction into Cassandra.
func (c *CassandraService) InsertTransaction(tx Transaction) error {
	return c.Session.Query(`
		INSERT INTO transactions (id, amount, type, status)
		VALUES (?, ?, ?, ?)
	`, tx.ID, tx.Amount, tx.Type, tx.Status).Exec()
}

// GetTransaction retrieves a transaction by ID from Cassandra.
func (c *CassandraService) GetTransaction(id string) (Transaction, error) {
	var tx Transaction
	err := c.Session.Query(`
		SELECT id, amount, type, status FROM transactions WHERE id = ?
	`, id).Consistency(gocql.One).Scan(&tx.ID, &tx.Amount, &tx.Type, &tx.Status)
	return tx, err
}