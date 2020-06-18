package awss3

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"io"
	"strings"
)

type SaveRequest struct {
	Path                string    `form:"path" json:"path"`
	File                io.Reader `form:"file" json:"file"`
	EncryptionKey       string    `form:"encryption_key" json:"encryption_key"`
	EncryptionAlgorithm string    `form:"encryption_algorithm" json:"encryption_algorithm"`
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

func (req *SaveRequest) UseEncryption() bool {
	if strings.TrimSpace(req.EncryptionKey) != "" {
		return true
	}
	return false
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
