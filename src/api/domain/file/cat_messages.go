package file

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"io"
)

type CatRequest struct {
	Path string `form:"path" json:"path"`
}

func (req *CatRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

type CatResponse struct {
	File        io.ReadCloser `form:"file" json:"file"`
	ContentType string        `form:"content_type" json:"content_type"`
}

func NewCatResponse(file io.ReadCloser, contentType string) *CatResponse {
	return &CatResponse{
		File:        file,
		ContentType: contentType,
	}
}
