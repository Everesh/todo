package models

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nirabyte/todo/internal/config"
	"github.com/nirabyte/todo/internal/storage"
)

var storageManager *storage.StorageManager

// InitStorage initializes the storage backend
func InitStorage(storageType string) error {
	var store storage.Storage
	var err error

	// Normalize storage type
	if storageType == "" {
		storageType = "file"
	}

	switch storageType {
	case "file":
		store, err = storage.NewFileStorage(config.DataPath)
		if err != nil {
			return fmt.Errorf("failed to initialize file storage: %w", err)
		}
	case "s3":
		if config.S3Bucket == "" {
			return fmt.Errorf("S3_BUCKET environment variable is required for S3 storage")
		}
		store, err = storage.NewS3Storage(config.S3Bucket, config.S3Region)
		if err != nil {
			return fmt.Errorf("failed to initialize S3 storage: %w", err)
		}
	case "mongodb", "mongo":
		store, err = storage.NewMongoStorage(config.MongoURI, config.MongoDB, config.MongoCollection)
		if err != nil {
			return fmt.Errorf("failed to initialize MongoDB storage: %w", err)
		}
	case "postgres", "postgresql", "pg":
		store, err = storage.NewDBStorage("postgres", config.PostgresDSN, config.PostgresTable)
		if err != nil {
			return fmt.Errorf("failed to initialize PostgreSQL storage: %w", err)
		}
	default:
		return fmt.Errorf("unsupported storage type: %s (supported: file, s3, mongodb, postgres)", storageType)
	}

	// Initialize encryption if key is provided
	var encryptor storage.Encryptor
	if config.EncryptionKey != "" {
		key, err := hex.DecodeString(config.EncryptionKey)
		if err != nil {
			return fmt.Errorf("invalid encryption key (must be 64 hex characters): %w", err)
		}
		if len(key) != 32 {
			return fmt.Errorf("encryption key must be 32 bytes (64 hex characters), got %d bytes", len(key))
		}
		encryptor, err = storage.NewAESEncryptor(key)
		if err != nil {
			return fmt.Errorf("failed to initialize encryptor: %w", err)
		}
	}

	storageManager = storage.NewStorageManager(store, encryptor)
	return nil
}

// GenerateEncryptionKey creates a new 32-byte key for AES-256
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func LoadData() AppData {
	hints := []Task{
		{ID: 1, Title: "Welcome to your TODO Manager", Done: false},
		{ID: 2, Title: "Press 'n' to add a new task", Done: false},
		{ID: 3, Title: "Press 'e' to edit the selected task", Done: false},
		{ID: 4, Title: "Press 'd' to delete a task", Done: false},
		{ID: 5, Title: "Press 'space' to check/uncheck", Done: true},
		{ID: 6, Title: "Press '@' to set a timer notification", Done: false},
		{ID: 7, Title: "Press 's' to cycle sort modes", Done: false},
		{ID: 8, Title: "Press 't' to change the color theme", Done: false},
	}

	defaultData := AppData{
		ThemeIndex: 0,
		SortMode:   SortOff,
		Tasks:      hints,
	}

	if storageManager == nil {
		return defaultData
	}

	data, err := storageManager.Load(config.DataFile)
	if err != nil {
		return defaultData
	}

	var appData AppData
	if err := json.Unmarshal(data, &appData); err == nil {
		for i := range appData.Tasks {
			if appData.Tasks[i].ID == 0 {
				appData.Tasks[i].ID = time.Now().UnixNano() + int64(i)
			}
		}
		return appData
	}
	return defaultData
}

func (m *Model) Save() {
	if storageManager == nil {
		// storageManager is nil we return here to avoid a panic
		return
	}

	var validTasks []Task
	for _, t := range m.Tasks {
		if !t.IsDeleting {
			validTasks = append(validTasks, t)
		}
	}

	data := AppData{
		ThemeIndex: m.ThemeIndex,
		SortMode:   m.SortMode,
		Tasks:      validTasks,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	_ = storageManager.Save(config.DataFile, bytes)
}
