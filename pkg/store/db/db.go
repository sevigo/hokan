package db

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

// DB holds the database connection
type DB struct {
	conn *bolt.DB
	path string
}

var ErrBucketNotFound = errors.New("bucket was not found")

// Add basic funtcions here
func (db *DB) Write(bucketName, key, value string) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func (db *DB) Read(bucketName, key string) ([]byte, error) {
	var result []byte
	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return ErrBucketNotFound
		}
		result = b.Get([]byte(key))
		return nil
	})
	return result, err
}

func (db *DB) ReadBucket(bucketName string) (map[string]string, error) {
	result := make(map[string]string)
	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return ErrBucketNotFound
		}
		return b.ForEach(func(k, v []byte) error {
			result[string(k)] = string(v)
			return nil
		})
	})
	return result, err
}
