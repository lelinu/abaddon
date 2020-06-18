package awss3

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"strings"
)

type TouchRequest struct {
	Path                string `form:"path" json:"path"`
	EncryptionKey       string `form:"encryption_key" json:"encryption_key"`
	EncryptionAlgorithm string `form:"encryption_algorithm" json:"encryption_algorithm"`
}

func (req *TouchRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

func (req *TouchRequest) UseEncryption() bool {
	if strings.TrimSpace(req.EncryptionKey) != "" {
		return true
	}
	return false
}
