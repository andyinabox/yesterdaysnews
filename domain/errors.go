package domain

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
	error
	Type() string
}

type ErrorHandler interface {
	Add(typ string, err error)
	Channel() chan<- Error
	Count(string) int
	CountAll() int
	Report() string
	DeferredReport()
	Reset()
	Err() error
}
