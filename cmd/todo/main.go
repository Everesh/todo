package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nirabyte/todo/internal/app"
	"github.com/nirabyte/todo/internal/config"
	"github.com/nirabyte/todo/internal/models"
)

func main() {
	// Setup logging to both file and stdout
	if err := setupLogging(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to setup logging: %v\n", err)
	}

	log.Println("=== TODO Application Starting ===")

	// Load configuration from environment variables
	log.Println("Loading configuration...")
	loadConfig()

	// Initialize or load encryption key
	log.Println("Initializing encryption key...")
	keyPath := getKeyPath()
	encryptionKey, isNewKey, err := initializeEncryptionKey(keyPath)
	if err != nil {
		log.Fatalf("Failed to initialize encryption key: %v", err)
	}

	// Set the encryption key in config
	config.EncryptionKey = encryptionKey

	// Log startup information
	logStartupInfo(keyPath, isNewKey)

	// Initialize storage backend with encryption
	log.Println("Initializing storage backend...")
	if err := models.InitStorage(config.StorageType); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	log.Println("Storage initialized successfully")

	// Load initial data
	log.Println("Loading application data...")

	// Create and run the application
	log.Println("Starting TUI application...")
	application := app.New()
	if err := application.Run(); err != nil {
		log.Printf("Application error: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	log.Println("Application exited normally")
}

func setupLogging() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(homeDir, ".todo")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logDir, "todo.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Write to both file and stdout
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Also print log location to stderr so user sees it
	fmt.Fprintf(os.Stderr, "Logging to: %s\n", logPath)

	return nil
}

func getKeyPath() string {
	// Check if user specified a custom key path
	if customPath := os.Getenv("TODO_KEY_PATH"); customPath != "" {
		log.Printf("Using custom key path from TODO_KEY_PATH: %s", customPath)
		return customPath
	}

	// Default to ~/.todo/key
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	keyPath := filepath.Join(homeDir, ".todo", "key")
	log.Printf("Using default key path: %s", keyPath)
	return keyPath
}

func initializeEncryptionKey(keyPath string) (string, bool, error) {
	log.Printf("Checking for encryption key at: %s", keyPath)

	// Check if key file exists
	if _, err := os.Stat(keyPath); err == nil {
		log.Println("Existing key file found")
		// Key exists, read it
		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return "", false, fmt.Errorf("failed to read existing key: %w", err)
		}

		key := strings.TrimSpace(string(keyBytes))

		// Validate key format
		if len(key) != 64 {
			return "", false, fmt.Errorf("invalid key length in %s: expected 64 hex characters, got %d", keyPath, len(key))
		}

		// Verify it's valid hex
		if _, err := hex.DecodeString(key); err != nil {
			return "", false, fmt.Errorf("invalid key format in %s: must be 64 hex characters", keyPath)
		}

		log.Println("Key validated successfully")
		return key, false, nil
	}

	// Key doesn't exist, generate a new one
	log.Println("No existing key found, generating new key...")
	key, err := generateEncryptionKey()
	if err != nil {
		return "", false, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(keyPath)
	log.Printf("Creating key directory: %s", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", false, fmt.Errorf("failed to create key directory: %w", err)
	}

	// Write the key to file with restrictive permissions
	log.Printf("Writing key to file: %s", keyPath)
	if err := os.WriteFile(keyPath, []byte(key), 0600); err != nil {
		return "", false, fmt.Errorf("failed to save encryption key: %w", err)
	}

	log.Println("New key generated and saved successfully")
	return key, true, nil
}

func generateEncryptionKey() (string, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func logStartupInfo(keyPath string, isNewKey bool) {
	log.Println("=== Startup Configuration ===")

	// Redact the key path to show only the directory structure
	displayPath := keyPath
	homeDir, err := os.UserHomeDir()
	if err == nil {
		displayPath = strings.Replace(keyPath, homeDir, "~", 1)
	}

	log.Printf("Storage type: %s", config.StorageType)
	log.Printf("Encryption key path: %s", displayPath)

	if isNewKey {
		log.Println("Encryption key: NEW (generated)")
	} else {
		log.Println("Encryption key: EXISTING (loaded)")
	}

	// Redacted key preview (show first 8 and last 8 characters)
	if config.EncryptionKey != "" {
		redacted := config.EncryptionKey[:8] + "..." + config.EncryptionKey[len(config.EncryptionKey)-8:]
		log.Printf("Encryption key (redacted): %s", redacted)
	}

	// Log storage-specific info
	switch config.StorageType {
	case "file":
		log.Printf("Data path: %s", config.DataPath)
		log.Printf("Data file: %s", config.DataFile)
		fullPath := filepath.Join(config.DataPath, config.DataFile)
		log.Printf("Full data path: %s", fullPath)
	case "s3":
		log.Printf("S3 bucket: %s", config.S3Bucket)
		log.Printf("S3 region: %s", config.S3Region)
	case "mongodb", "mongo":
		log.Printf("MongoDB URI: %s", redactMongoURI(config.MongoURI))
		log.Printf("MongoDB database: %s", config.MongoDB)
		log.Printf("MongoDB collection: %s", config.MongoCollection)
	case "postgres", "postgresql", "pg":
		// Redact password from DSN
		dsn := redactDSN(config.PostgresDSN)
		log.Printf("PostgreSQL DSN: %s", dsn)
		log.Printf("PostgreSQL table: %s", config.PostgresTable)
	}

	log.Println("=============================")
}

func redactDSN(dsn string) string {
	// Simple redaction for PostgreSQL DSN
	// postgres://user:password@host/db -> postgres://user:***@host/db
	if strings.Contains(dsn, "://") && strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "://")
		if len(parts) == 2 {
			userPart := strings.Split(parts[1], "@")
			if len(userPart) == 2 {
				credentials := strings.Split(userPart[0], ":")
				if len(credentials) == 2 {
					return parts[0] + "://" + credentials[0] + ":***@" + userPart[1]
				}
			}
		}
	}
	return dsn
}

func redactMongoURI(uri string) string {
	// mongodb://user:password@host/db -> mongodb://user:***@host/db
	if strings.Contains(uri, "://") && strings.Contains(uri, "@") {
		parts := strings.Split(uri, "://")
		if len(parts) == 2 {
			userPart := strings.Split(parts[1], "@")
			if len(userPart) == 2 {
				credentials := strings.Split(userPart[0], ":")
				if len(credentials) == 2 {
					return parts[0] + "://" + credentials[0] + ":***@" + userPart[1]
				}
			}
		}
	}
	return uri
}

func loadConfig() {
	// Storage type (defaults to "file" if not set)
	if storageType := os.Getenv("STORAGE_TYPE"); storageType != "" {
		config.StorageType = storageType
		log.Printf("Config: STORAGE_TYPE=%s", storageType)
	}

	// S3 configuration
	if s3Bucket := os.Getenv("S3_BUCKET"); s3Bucket != "" {
		config.S3Bucket = s3Bucket
		log.Printf("Config: S3_BUCKET=%s", s3Bucket)
	}
	if s3Region := os.Getenv("S3_REGION"); s3Region != "" {
		config.S3Region = s3Region
		log.Printf("Config: S3_REGION=%s", s3Region)
	}

	// MongoDB configuration
	if mongoURI := os.Getenv("MONGO_URI"); mongoURI != "" {
		config.MongoURI = mongoURI
		log.Printf("Config: MONGO_URI=%s", redactMongoURI(mongoURI))
	}
	if mongoDB := os.Getenv("MONGO_DB"); mongoDB != "" {
		config.MongoDB = mongoDB
		log.Printf("Config: MONGO_DB=%s", mongoDB)
	}
	if mongoCollection := os.Getenv("MONGO_COLLECTION"); mongoCollection != "" {
		config.MongoCollection = mongoCollection
		log.Printf("Config: MONGO_COLLECTION=%s", mongoCollection)
	}

	// PostgreSQL configuration
	if pgDSN := os.Getenv("POSTGRES_DSN"); pgDSN != "" {
		config.PostgresDSN = pgDSN
		log.Printf("Config: POSTGRES_DSN=%s", redactDSN(pgDSN))
	}
	if pgTable := os.Getenv("POSTGRES_TABLE"); pgTable != "" {
		config.PostgresTable = pgTable
		log.Printf("Config: POSTGRES_TABLE=%s", pgTable)
	}

	// Data file configuration
	if dataPath := os.Getenv("DATA_PATH"); dataPath != "" {
		config.DataPath = dataPath
		log.Printf("Config: DATA_PATH=%s", dataPath)
	}
	if dataFile := os.Getenv("DATA_FILE"); dataFile != "" {
		config.DataFile = dataFile
		log.Printf("Config: DATA_FILE=%s", dataFile)
	}

	log.Println("Configuration loaded successfully")
}
