// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"crypto-exchange/config"
	"crypto-exchange/controllers"
	"crypto-exchange/middleware"
	"crypto-exchange/routes"
	"crypto-exchange/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	// Load configuration from config.yaml
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := cfg.SetupLogger()

	// Initialize database service
	dbService, err := services.NewDatabaseService(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database service")
	}

	// Initialize Redis service
	redisService := services.NewRedisService(cfg.Redis)
	logger.Info().Msg("Connected to Redis")

	// Initialize Cassandra service
	cassandraService, err := services.NewCassandraService(cfg.Cassandra)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize Cassandra service")
	}
	defer cassandraService.Close()
	logger.Info().Msg("Connected to Cassandra")

	// Initialize Kafka service
	kafkaService := services.NewKafkaService(cfg.Kafka)
	defer kafkaService.Close()
	logger.Info().Msg("Connected to Kafka")

	// Choose between mock service or real database service based on environment
	var txService services.TransactionService
	if cfg.Environment == "development" {
		txService = services.NewMockTransactionService()
		logger.Info().Msg("Using MockTransactionService")
	} else {
		// Initialize production transaction service with DB, Redis, Cassandra, Kafka
		txService = services.NewTransactionService(dbService.DB, logger, redisService, cassandraService, kafkaService)
		logger.Info().Msg("Using TransactionService with PostgreSQL, Redis, Cassandra, Kafka")
	}

	// Initialize controllers
	txController := controllers.NewTransactionController(txService, logger)

	// Initialize Gin router
	router := gin.New()

	// Apply middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))

	// Setup routes
	routes.SetupRoutes(router, txController, logger)

	// Configure server settings
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           serverAddr,
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	logger.Info().
		Str("address", serverAddr).
		Msg("Starting server")

	// Start server
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().
			Err(err).
			Msg("Server failed")
	}

	// Implement graceful shutdown if needed
}