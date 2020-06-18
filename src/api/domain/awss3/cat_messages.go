package awss3

import (
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"io"
	"strings"
)

type CatRequest struct {
	Path                string `form:"path" json:"path"`
	EncryptionKey       string `form:"encryption_key" json:"encryption_key"`
	EncryptionAlgorithm string `form:"encryption_algorithm" json:"encryption_algorithm"`
}

func (req *CatRequest) Validate() *error_utils.ApiError {

	// set defaults
	req.setEncryptionAlgorithm()

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

func (req *CatRequest) UseEncryption() bool {
	if strings.TrimSpace(req.EncryptionKey) != "" {
		return true
	}
	return false
}

type CatResponse struct {
	File io.ReadCloser `form:"file" json:"file"`
}

func NewCatResponse(file io.ReadCloser) *CatResponse {
	return &CatResponse{File: file}
}

func (req *CatRequest) setEncryptionAlgorithm() {
	if strings.TrimSpace(req.EncryptionKey) != "" && strings.TrimSpace(req.EncryptionAlgorithm) != "" {
		req.EncryptionAlgorithm = config.GetS3EncryptionAlgorithm()
	}
}
