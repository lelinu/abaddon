package file

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
)

type MvRequest struct {
	PathFrom            string `form:"path_from" json:"path_from"`
	PathTo              string `form:"path_to" json:"path_to"`
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


