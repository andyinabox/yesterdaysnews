package domain

import (
	"encoding/json"
	"fmt"
)

const (
	// larger build phases
	ErrTypeGetVideoID    = "ErrTypeGetVideoID"
	ErrTypeDownloadVideo = "ErrTypeDownloadVideo"
	ErrTypeCutVideo      = "ErrTypeCutVideo"
	ErrTypeCleanup       = "ErrTypeCleanup"
	ErrTypeExtractImage  = "ErrTypeExtractImage"
	ErrTypeAverageImage  = "ErrTypeAverageImage"

	// object store
	ErrTypeCopyObject   = "ErrTypeCopyObject"
	ErrTypeMoveObject   = "ErrTypeMoveObject"
	ErrTypeDeleteObject = "ErrTypeDeleteObject"
	ErrTypeUploadFile   = "ErrTypeUploadFile"

	// misc
	ErrTypeFatal = "ErrTypeFatal"
	ErrTypeTODO  = "ErrTypeTODO"
)

type Error interface {
	error
	json.Marshaler
	fmt.Stringer
	Type() string
}

type ErrorHandler interface {
	fmt.Stringer
	json.Marshaler

	// adding errors
	Add(typ string, err error)
	Channel() chan<- Error

	// error counts
	Count(string) int
	CountAll() int

	// output
	Report()
}
