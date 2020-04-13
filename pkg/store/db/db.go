package db

import (
	"fmt"

	"github.com/sevigo/hokan/pkg/core"
	bolt "go.etcd.io/bbolt"
)

// Bolt holds the database connection
type Bolt struct {
	conn *bolt.DB
	path string
}

// Add basic funtcions here
func (db *Bolt) Write(bucketName, key, value string) error {
	return db.conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func (db *Bolt) Read(bucketName, key string) ([]byte, error) {
	var result []byte
	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return &ErrBucketNotFound{fmt.Sprintf("Read: bucket %q was not found", bucketName)}
		}
		result = b.Get([]byte(key))
		return nil
	})
	return result, err
}

func (db *Bolt) ReadBucket(bucketName string, opt *core.ReadBucketOptions) (map[string]string, error) {
	result := make(map[string]string)
	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return &ErrBucketNotFound{fmt.Sprintf("ReadBucket: bucket %q was not found", bucketName)}
		}
		return b.ForEach(func(k, v []byte) error {
			result[string(k)] = string(v)
			return nil
		})
		// TODO: use Cursor
		// c := b.Cursor()
		// for k, v := c.First(); k != nil; k, v = c.Next() {
		// 	fmt.Printf("key=%s, value=%s\n", k, v)
		// }
	})
	return result, err
}
