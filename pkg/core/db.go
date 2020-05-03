package core

type ReadBucketOptions struct {
	Query  string
	Offset uint64
	Limit  uint64
}

type BucketData struct {
	Key   string
	Value string
}

// DB interface for the bold storage
type DB interface {
	Write(bucketName, key, value string) error
	Read(bucketName, key string) ([]byte, error)
	ReadBucket(bucketName string, opt *ReadBucketOptions) ([]BucketData, error)
}
