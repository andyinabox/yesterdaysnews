package domain

import "gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"

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
