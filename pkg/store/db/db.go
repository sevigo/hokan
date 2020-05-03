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

func (db *Bolt) ReadBucket(bucketName string, opt *core.ReadBucketOptions) ([]core.BucketData, error) {
	if opt == nil {
		opt = &core.ReadBucketOptions{}
	}
	result := []core.BucketData{}

	err := db.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return &ErrBucketNotFound{fmt.Sprintf("ReadBucket: bucket %q was not found", bucketName)}
		}
		c := b.Cursor()
		i := uint64(0)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i++
			if opt.Limit > 0 {
				if i-1 < opt.Offset {
					continue
				}
				if uint64(len(result)) < opt.Limit {
					result = append(result, core.BucketData{
						Key:   string(k),
						Value: string(v),
					})
				} else {
					break
				}
			} else {
				result = append(result, core.BucketData{
					Key:   string(k),
					Value: string(v),
				})
			}
		}
		return nil
	})
	return result, err
}
