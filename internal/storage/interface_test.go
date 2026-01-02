package storage

import (
	"bytes"
	"errors"
	"testing"
)

type mockStorage struct {
	data      map[string][]byte
	loadErr   error
	saveErr   error
	deleteErr error
	existsErr error
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		data: make(map[string][]byte),
	}
}

func (m *mockStorage) Load(key string) ([]byte, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	v, ok := m.data[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *mockStorage) Save(key string, data []byte) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.data[key] = data
	return nil
}

func (m *mockStorage) Delete(key string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.data, key)
	return nil
}

func (m *mockStorage) Exists(key string) (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}
	_, ok := m.data[key]
	return ok, nil
}

func TestAESEncryptor_New(t *testing.T) {
	key := make([]byte, 32)
	_, err := NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = NewAESEncryptor([]byte("short"))
	if err == nil {
		t.Fatalf("expected error for invalid key length")
	}
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	enc, err := NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plaintext := []byte("secret-data")
	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	result, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(result, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, result)
	}
}

func TestAESEncryptor_Decrypt_InvalidCiphertext(t *testing.T) {
	key := make([]byte, 32)
	enc, _ := NewAESEncryptor(key)

	_, err := enc.Decrypt([]byte("short"))
	if err == nil {
		t.Fatalf("expected error for short ciphertext")
	}
}

func TestStorageManager_SaveLoad_NoEncryption(t *testing.T) {
	st := newMockStorage()
	sm := NewStorageManager(st, nil)

	key := "k1"
	value := []byte("plain")

	if err := sm.Save(key, value); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	out, err := sm.Load(key)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !bytes.Equal(out, value) {
		t.Fatalf("expected %q, got %q", value, out)
	}
}

func TestStorageManager_SaveLoad_WithEncryption(t *testing.T) {
	st := newMockStorage()
	keyBytes := make([]byte, 32)
	enc, _ := NewAESEncryptor(keyBytes)
	sm := NewStorageManager(st, enc)

	key := "k1"
	value := []byte("encrypted")

	if err := sm.Save(key, value); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if bytes.Equal(st.data[key], value) {
		t.Fatalf("data should be encrypted at rest")
	}

	out, err := sm.Load(key)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !bytes.Equal(out, value) {
		t.Fatalf("expected %q, got %q", value, out)
	}
}

func TestStorageManager_Load_Error(t *testing.T) {
	st := newMockStorage()
	st.loadErr = errors.New("load error")
	sm := NewStorageManager(st, nil)

	_, err := sm.Load("k")
	if err == nil {
		t.Fatalf("expected load error")
	}
}

func TestStorageManager_Save_EncryptError(t *testing.T) {
	st := newMockStorage()
	enc := &AESEncryptor{key: []byte("bad")}
	sm := NewStorageManager(st, enc)

	err := sm.Save("k", []byte("data"))
	if err == nil {
		t.Fatalf("expected encryption error")
	}
}

func TestStorageManager_Delete(t *testing.T) {
	st := newMockStorage()
	sm := NewStorageManager(st, nil)

	st.data["k"] = []byte("v")
	if err := sm.Delete("k"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, ok := st.data["k"]; ok {
		t.Fatalf("expected key to be deleted")
	}
}

func TestStorageManager_Exists(t *testing.T) {
	st := newMockStorage()
	sm := NewStorageManager(st, nil)

	exists, err := sm.Exists("k")
	if err != nil {
		t.Fatalf("exists failed: %v", err)
	}
	if exists {
		t.Fatalf("expected key to not exist")
	}

	st.data["k"] = []byte("v")
	exists, _ = sm.Exists("k")
	if !exists {
		t.Fatalf("expected key to exist")
	}
}
