package db

import bolt "go.etcd.io/bbolt"

// DB holds the database connection
type DB struct {
	conn *bolt.DB
	path string
}

// Add basic funtcions here
