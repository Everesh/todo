package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileStorage(t *testing.T) {
	dir := t.TempDir()

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fs.basePath != dir {
		t.Fatalf("unexpected basePath")
	}
}

func TestFileStorage_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	key := "test-key"
	data := []byte("hello world")

	if err := fs.Save(key, data); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	out, err := fs.Load(key)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !bytes.Equal(out, data) {
		t.Fatalf("expected %q, got %q", data, out)
	}
}

func TestFileStorage_Load_NotFound(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	_, err := fs.Load("missing")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestFileStorage_Save_Error(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	readOnlyDir := filepath.Join(dir, "ro")
	if err := os.Mkdir(readOnlyDir, 0555); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	fs.basePath = readOnlyDir

	err := fs.Save("key", []byte("data"))
	if err == nil {
		t.Fatalf("expected write error")
	}
}

func TestFileStorage_Delete(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	key := "k"
	if err := fs.Save(key, []byte("v")); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := fs.Delete(key); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := fs.Load(key); err == nil {
		t.Fatalf("expected error after delete")
	}
}

func TestFileStorage_Delete_NotFound(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	err := fs.Delete("missing")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestFileStorage_Exists(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	exists, err := fs.Exists("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatalf("expected not exists")
	}

	if err := fs.Save("k", []byte("v")); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	exists, err = fs.Exists("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatalf("expected exists")
	}
}

func TestFileStorage_Exists_Error(t *testing.T) {
	dir := t.TempDir()
	fs, _ := NewFileStorage(dir)

	badPath := filepath.Join(dir, "file")
	if err := os.WriteFile(badPath, []byte("x"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	fs.basePath = badPath

	_, err := fs.Exists("k")
	if err == nil {
		t.Fatalf("expected stat error")
	}
}
