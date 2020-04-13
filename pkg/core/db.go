package core

type ReadBucketOptions struct {
	Query  string
	Offset int64
	Limit  int64
}

// DB interface for the bold storage
type DB interface {
	Write(bucketName, key, value string) error
	Read(bucketName, key string) ([]byte, error)
	ReadBucket(bucketName string, opt *ReadBucketOptions) (map[string]string, error)
}
