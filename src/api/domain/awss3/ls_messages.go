package awss3

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"os"
)

type LsRequest struct {
	Path string `form:"path" json:"path"`
}

func (req *LsRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

type LsResponse struct {
	Files []os.FileInfo `form:"files" json:"files"`
}

func NewLsResponse(files []os.FileInfo) *LsResponse {
	return &LsResponse{Files: files}
}
