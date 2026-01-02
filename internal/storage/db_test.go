package storage

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockDBStorage(t *testing.T) (*DBStorage, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	storage := &DBStorage{
		db:        db,
		tableName: "test_table",
	}

	return storage, mock, func() {
		db.Close()
	}
}

func TestDBStorage_Save(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO test_table (key, data) 
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET data = $2`,
	)).
		WithArgs("k1", []byte("v1")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := ds.Save("k1", []byte("v1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBStorage_Save_Error(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO test_table (key, data) 
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET data = $2`,
	)).
		WillReturnError(errors.New("insert error"))

	err := ds.Save("k", []byte("v"))
	if err == nil {
		t.Fatalf("expected save error")
	}
}

func TestDBStorage_Load_Success(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"data"}).
		AddRow([]byte("value"))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT data FROM test_table WHERE key = $1",
	)).
		WithArgs("k1").
		WillReturnRows(rows)

	data, err := ds.Load("k1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "value" {
		t.Fatalf("unexpected data: %s", data)
	}
}

func TestDBStorage_Load_NotFound(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT data FROM test_table WHERE key = $1",
	)).
		WithArgs("k").
		WillReturnError(sql.ErrNoRows)

	_, err := ds.Load("k")
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestDBStorage_Load_Error(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT data FROM test_table WHERE key = $1",
	)).
		WillReturnError(errors.New("query error"))

	_, err := ds.Load("k")
	if err == nil {
		t.Fatalf("expected query error")
	}
}

func TestDBStorage_Delete(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM test_table WHERE key = $1",
	)).
		WithArgs("k").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := ds.Delete("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDBStorage_Exists(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS(SELECT 1 FROM test_table WHERE key = $1)",
	)).
		WithArgs("k").
		WillReturnRows(rows)

	exists, err := ds.Exists("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatalf("expected key to exist")
	}
}

func TestDBStorage_Exists_Error(t *testing.T) {
	ds, mock, cleanup := newMockDBStorage(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS(SELECT 1 FROM test_table WHERE key = $1)",
	)).
		WillReturnError(errors.New("exists error"))

	_, err := ds.Exists("k")
	if err == nil {
		t.Fatalf("expected exists error")
	}
}
