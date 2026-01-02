package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Storage interface for different storage backends
type Storage interface {
	Load(key string) ([]byte, error)
	Save(key string, data []byte) error
	Delete(key string) error
	Exists(key string) (bool, error)
}

// Encryptor handles data encryption/decryption
type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// AESEncryptor implements AES-GCM encryption
type AESEncryptor struct {
	key []byte
}

func NewAESEncryptor(key []byte) (*AESEncryptor, error) {
	if len(key) != 32 { // AES-256
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	return &AESEncryptor{key: key}, nil
}

func (e *AESEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (e *AESEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// StorageManager
type StorageManager struct {
	storage   Storage
	encryptor Encryptor
}

func NewStorageManager(storage Storage, encryptor Encryptor) *StorageManager {
	return &StorageManager{
		storage:   storage,
		encryptor: encryptor,
	}
}

func (sm *StorageManager) Load(key string) ([]byte, error) {
	data, err := sm.storage.Load(key)
	if err != nil {
		return nil, err
	}

	if sm.encryptor != nil {
		return sm.encryptor.Decrypt(data)
	}

	return data, nil
}

func (sm *StorageManager) Save(key string, data []byte) error {
	var toSave []byte
	var err error

	if sm.encryptor != nil {
		toSave, err = sm.encryptor.Encrypt(data)
		if err != nil {
			return err
		}
	} else {
		toSave = data
	}

	return sm.storage.Save(key, toSave)
}

func (sm *StorageManager) Delete(key string) error {
	return sm.storage.Delete(key)
}

func (sm *StorageManager) Exists(key string) (bool, error) {
	return sm.storage.Exists(key)
}
