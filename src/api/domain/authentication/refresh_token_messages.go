package authentication

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
)

type RefreshTokenRequest struct {
	Token     string `json:"token"`
}

func (req *RefreshTokenRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()
	v.IsNotEmpty("token", req.Token)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

type RefreshTokenResponse struct {
	AuthenticationToken string `json:"token"`
	ExpiresAt           string `json:"expires_at"`
}

func NewRefreshTokenResponse(authenticationToken string, expiresAt string) *RefreshTokenResponse {
	return &RefreshTokenResponse{AuthenticationToken: authenticationToken, ExpiresAt: expiresAt}
}
