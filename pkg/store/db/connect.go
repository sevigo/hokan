package db

import (
	"time"

	"github.com/sevigo/hokan/pkg/core"
	bolt "go.etcd.io/bbolt"
)

// Connect open the bolt database for r/w
func Connect(path string) (core.DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Bolt{
		conn: db,
		path: path,
	}, nil
}

// Close the database connection
func (db *Bolt) Close() error {
	return db.conn.Close()
}
