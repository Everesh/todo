package storage

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBStorage implements SQL database storage
type DBStorage struct {
	db        *sql.DB
	tableName string
}

func NewDBStorage(driverName, dataSourceName, tableName string) (*DBStorage, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	storage := &DBStorage{
		db:        db,
		tableName: tableName,
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			key VARCHAR(255) PRIMARY KEY,
			data BYTEA NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (ds *DBStorage) Load(key string) ([]byte, error) {
	var data []byte
	err := ds.db.QueryRow(
		"SELECT data FROM "+ds.tableName+" WHERE key = $1",
		key,
	).Scan(&data)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("key not found")
		}
		return nil, err
	}

	return data, nil
}

func (ds *DBStorage) Save(key string, data []byte) error {
	_, err := ds.db.Exec(`
		INSERT INTO `+ds.tableName+` (key, data) 
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET data = $2
	`, key, data)
	return err
}

func (ds *DBStorage) Delete(key string) error {
	_, err := ds.db.Exec("DELETE FROM "+ds.tableName+" WHERE key = $1", key)
	return err
}

func (ds *DBStorage) Exists(key string) (bool, error) {
	var exists bool
	err := ds.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+ds.tableName+" WHERE key = $1)",
		key,
	).Scan(&exists)
	return exists, err
}

func (ds *DBStorage) Close() error {
	return ds.db.Close()
}
