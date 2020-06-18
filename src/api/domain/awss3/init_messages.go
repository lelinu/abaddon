package awss3

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
)

type InitRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
	Region   string `json:"region"`
}

func NewInitRequest(username string, password string, path string, region string) *InitRequest {

	return &InitRequest{
		Username: username,
		Password: password,
		Path:     path,
		Region:   region,
	}
}

func (req *InitRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("username", req.Username)
	v.IsNotEmpty("password", req.Password)
	v.IsNotEmpty("path", req.Path)
	v.IsNotEmpty("region", req.Region)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}
