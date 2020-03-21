package db

type ErrBucketNotFound struct {
	msg string
}

func (e *ErrBucketNotFound) Error() string { return e.msg }
