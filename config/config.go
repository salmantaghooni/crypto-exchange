// config/config.go
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Config represents the entire configuration structure.
type Config struct {
	Environment      string                 `mapstructure:"environment" validate:"required,oneof=development production"`
	Server           ServerConfig           `mapstructure:"server" validate:"required,dive"`
	Database         DatabaseConfig         `mapstructure:"database" validate:"required,dive"`
	Logging          LoggingConfig          `mapstructure:"logging" validate:"required,dive"`
	JWT              JWTConfig              `mapstructure:"jwt" validate:"required,dive"`
	APIKeys          APIKeysConfig          `mapstructure:"api_keys" validate:"required,dive"`
	ExternalServices ExternalServicesConfig `mapstructure:"external_services" validate:"required,dive"`
	Redis            RedisConfig            `mapstructure:"redis" validate:"required,dive"`
	Cassandra        CassandraConfig        `mapstructure:"cassandra" validate:"required,dive"`
	Kafka            KafkaConfig            `mapstructure:"kafka" validate:"required,dive"`
	Features         FeaturesConfig         `mapstructure:"features"`
}

// ServerConfig holds server-related configurations.
type ServerConfig struct {
	Host         string        `mapstructure:"host" validate:"required,ip"`
	Port         int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" validate:"required"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" validate:"required"`
}

// DatabaseConfig holds database-related configurations.
type DatabaseConfig struct {
	Type     string         `mapstructure:"type" validate:"required,oneof=postgres mysql sqlite"`
	Postgres PostgresConfig `mapstructure:"postgres" validate:"required_if=Type postgres"`
	MySQL    MySQLConfig    `mapstructure:"mysql" validate:"required_if=Type mysql"`
}

// PostgresConfig holds PostgreSQL-specific configurations.
type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" validate:"required"`
}

// MySQLConfig holds MySQL-specific configurations.
type MySQLConfig struct {
	Host     string `mapstructure:"host" validate:"required_if=Type mysql"`
	Port     int    `mapstructure:"port" validate:"required_if=Type mysql"`
	User     string `mapstructure:"user" validate:"required_if=Type mysql"`
	Password string `mapstructure:"password" validate:"required_if=Type mysql"`
	DBName   string `mapstructure:"dbname" validate:"required_if=Type mysql"`
}

// LoggingConfig holds logging-related configurations.
type LoggingConfig struct {
	Level       string   `mapstructure:"level" validate:"required,oneof=debug info warn error fatal"`
	Format      string   `mapstructure:"format" validate:"required,oneof=console json"`
	OutputPaths []string `mapstructure:"output_paths" validate:"required,min=1,dive,required"`
}

// JWTConfig holds JWT-related configurations.
type JWTConfig struct {
	SecretKey     string        `mapstructure:"secret_key" validate:"required"`
	TokenDuration time.Duration `mapstructure:"token_duration" validate:"required"`
}

// APIKeysConfig holds API keys for external services.
type APIKeysConfig struct {
	CryptoAPI CryptoAPIKeys `mapstructure:"crypto_api" validate:"required,dive"`
}

// CryptoAPIKeys holds specific API keys for crypto services.
type CryptoAPIKeys struct {
	Key    string `mapstructure:"key" validate:"required"`
	Secret string `mapstructure:"secret" validate:"required"`
}

// ExternalServicesConfig holds configurations for external services.
type ExternalServicesConfig struct {
	PaymentGateway      ServiceConfig `mapstructure:"payment_gateway" validate:"required,dive"`
	ExchangeRateService ServiceConfig `mapstructure:"exchange_rate_service" validate:"required,dive"`
}

// ServiceConfig holds configurations for a generic service.
type ServiceConfig struct {
	BaseURL string `mapstructure:"base_url" validate:"required,url"`
	APIKey  string `mapstructure:"api_key" validate:"required"`
}

// RedisConfig holds Redis-related configurations.
type RedisConfig struct {
	Host     string `mapstructure:"host" validate:"required,ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db" validate:"required"`
}

// CassandraConfig holds Cassandra-related configurations.
type CassandraConfig struct {
	Host     string `mapstructure:"host" validate:"required,ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Keyspace string `mapstructure:"keyspace" validate:"required"`
}

// KafkaConfig holds Kafka-related configurations.
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers" validate:"required,min=1,dive,required"`
	Topic   string   `mapstructure:"topic" validate:"required"`
}

// FeaturesConfig holds feature flags configurations.
type FeaturesConfig struct {
	EnableNewFeatureX bool `mapstructure:"enable_new_feature_x"`
	EnableLogging     bool `mapstructure:"enable_logging"`
}

// LoadConfig reads configuration from config.yaml and environment variables.
func LoadConfig(configPath string) (Config, error) {
	var config Config

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "console")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the config into the Config struct
	err := viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Validate the configuration
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return config, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// SetupLogger initializes the logger based on configuration.
func (c Config) SetupLogger() zerolog.Logger {
	var writer zerolog.Writer

	if c.Logging.Format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		writer = os.Stdout // JSON format
	}

	// Create a multi-writer for all output paths
	outputWriters := []zerolog.Writer{}
	for _, path := range c.Logging.OutputPaths {
		if path == "stdout" {
			outputWriters = append(outputWriters, writer)
		} else {
			// Ensure the logs directory exists
			dir := "logs"
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.Mkdir(dir, 0755)
			}
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Could not open log file %s: %v\n", path, err)
				continue
			}
			outputWriters = append(outputWriters, file)
		}
	}

	multiWriter := zerolog.MultiLevelWriter(outputWriters...)
	logger := zerolog.New(multiWriter).With().
		Timestamp().
		Caller().
		Logger()

	// Set the log level
	level, err := zerolog.ParseLevel(c.Logging.Level)
	if err != nil {
		logger.Warn().Msgf("Invalid log level %s, defaulting to info", c.Logging.Level)
		level = zerolog.InfoLevel
	}
	logger = logger.Level(level)

	return logger
}