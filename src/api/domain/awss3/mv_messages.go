package awss3

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"strings"
)

type MvRequest struct {
	PathFrom            string `form:"path_from" json:"path_from"`
	PathTo              string `form:"path_to" json:"path_to"`
	EncryptionKey       string `form:"encryption_key" json:"encryption_key"`
	EncryptionAlgorithm string `form:"encryption_algorithm" json:"encryption_algorithm"`
}

func (req *MvRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path_from", req.PathFrom)
	v.IsNotEmpty("path_to", req.PathTo)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

func (req *MvRequest) UseEncryption() bool {
	if strings.TrimSpace(req.EncryptionKey) != "" {
		return true
	}
	return false
}
