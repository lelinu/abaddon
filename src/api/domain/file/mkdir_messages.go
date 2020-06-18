package file

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
)

type MkDirRequest struct {
	Path string `form:"path" json:"path" example:"/hello-test/"`
}

func (req *MkDirRequest) Validate() *error_utils.ApiError {
	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}
