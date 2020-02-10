package db

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

// Connect open the bolt database for r/w
func Connect(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &DB{
		conn: db,
		path: path,
	}, nil
}

// Close the database conection
func (db *DB) Close() error {
	return db.Close()
}
