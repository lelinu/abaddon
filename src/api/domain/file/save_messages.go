package file

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"mime/multipart"
)

type SaveRequest struct {
	Path string                `form:"path" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (req *SaveRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)
	v.IsNotNil("file", req.File)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

type SaveResponse struct {
	Location  string  `form:"location" json:"location"`
	VersionID *string `form:"version_id" json:"version_id"`
	UploadID  string  `form:"upload_id" json:"upload_id"`
}

func NewSaveResponse(location string, versionID *string, uploadID string) *SaveResponse {
	return &SaveResponse{
		Location:  location,
		VersionID: versionID,
		UploadID:  uploadID,
	}
}
