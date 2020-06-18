package awss3

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
)

type RmRequest struct {
	Path string `form:"path" json:"path"`
}

func (req *RmRequest) Validate() *error_utils.ApiError {
	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

func NewRmRequest(path string) *RmRequest {
	return &RmRequest{Path: path}
}
