package objectstoreclient

import "errors"

var (
	ErrBucketDoesNotExist = errors.New("bucket does not exist")
	ErrObjectDoesNotExist = errors.New("object does not exist")
)
