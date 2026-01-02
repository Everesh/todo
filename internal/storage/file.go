package storage

import (
	"os"
	"path/filepath"
)

// FileStorage implements local file storage
type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &FileStorage{basePath: basePath}, nil
}

func (fs *FileStorage) Load(key string) ([]byte, error) {
	path := filepath.Join(fs.basePath, key)
	return os.ReadFile(path)
}

func (fs *FileStorage) Save(key string, data []byte) error {
	path := filepath.Join(fs.basePath, key)
	return os.WriteFile(path, data, 0644)
}

func (fs *FileStorage) Delete(key string) error {
	path := filepath.Join(fs.basePath, key)
	return os.Remove(path)
}

func (fs *FileStorage) Exists(key string) (bool, error) {
	path := filepath.Join(fs.basePath, key)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
