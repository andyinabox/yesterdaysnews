package domain

import (
	"context"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"
)

const (
	// larger build phases
	ErrTypeGetVideoID               = "ErrTypeGetVideoID"
	ErrTypeDownloadVideo            = "ErrTypeDownloadVideo"
	ErrTypeCutVideo                 = "ErrTypeCutVideo"
	ErrTypeGetVideoEditPoints       = "ErrTypeGetVideoEditPoints"
	ErrTypeUploadVideo              = "ErrTypeUploadVideo"
	ErrTypeBuildModel               = "ErrTypeBuildModel"
	ErrTypeSaveModel                = "ErrTypeSaveModel"
	ErrTypeUploadModel              = "ErrTypeUploadModel"
	ErrTypeSaveManifest             = "ErrTypeSaveManifest"
	ErrTypeUploadManifest           = "ErrTypeUploadManifest"
	ErrTypeGetCurrentManifestPrefix = "ErrTypeGetCurrentManifestPrefix"
	ErrTypeCleanup                  = "ErrTypeCleanup"

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
	errorhandler.Error
}

type ErrorHandler interface {
	errorhandler.ErrorHandler
}

func Err(typ string, err error) Error {
	return errorhandler.Err(typ, err)
}

func DefaultErrorHandler(ctx context.Context) ErrorHandler {
	// set up error handling
	fatalFunc := func(typ string, err error) {
		log.Fatalf("%s: %s", typ, err)
	}
	errorFunc := func(typ string, err error) {
		if typ == ErrTypeFatal {
			fatalFunc(typ, err)
		}
		log.Errorf("%s: %s", typ, err)
	}
	return errorhandler.New(ctx, &errorhandler.Config{
		ErrorFunc: errorFunc,
		FatalFunc: fatalFunc,
		Thresholds: map[string]int{
			ErrTypeMoveObject: 10,
		},
	})
}
