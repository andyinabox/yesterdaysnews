package domain

import "gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"

const (
	ErrTypeTODO               = "ErrTypeTODO"
	ErrTypeGetVideoID         = "ErrTypeGetVideoID"
	ErrTypeDownloadVideo      = "ErrTypeDownloadVideo"
	ErrTypeCutVideo           = "ErrTypeCutVideo"
	ErrTypeGetVideoEditPoints = "ErrTypeGetVideoEditPoints"
	ErrTypeUploadVideo        = "ErrTypeUploadVideo"
	ErrTypeBuildModel         = "ErrTypeBuildModel"
	ErrTypeSaveModel          = "ErrTypeSaveModel"
	ErrTypeUploadModel        = "ErrTypeUploadModel"
	ErrTypeSaveManifest       = "ErrTypeSaveManifest"
	ErrTypeUploadManifest     = "ErrTypeUploadManifest"
	ErrTypeFatal              = "ErrTypeFatal"
)

type Error interface {
	errorhandler.Error
}

type ErrorHandler interface {
	errorhandler.ErrorHandler
}
