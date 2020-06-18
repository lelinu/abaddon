package file

import (
	"github.com/lelinu/abaddon/src/api/utils/file_utils"
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
)

type LsRequest struct {
	Path string `form:"path" json:"path" example:"/"`
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
	Files []file_utils.FileInfo `form:"files" json:"files"`
	Etag  string                `form:"etag" json:"etag"`
}

func NewLsResponse(files []file_utils.FileInfo, etag string) *LsResponse {
	return &LsResponse{Files: files, Etag: etag}
}
