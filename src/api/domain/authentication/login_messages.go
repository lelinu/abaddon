package authentication

import (
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/abaddon/src/api/utils/path_utils"
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"strings"
)

type LoginRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	EncryptionKey   string `json:"encryption_key"`
	Path            string `json:"path"`
	Region          string `json:"region"`
}

func (req *LoginRequest) Validate() *error_utils.ApiError {

	// set defaults
	req.setDefaultRegion()
	req.enforcePath()

	v := validator_utils.NewValidator()

	v.IsNotEmpty("username", req.Username)
	v.IsNotEmpty("password", req.Password)

	v.IsNotEmpty("region", req.Region)
	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

type LoginResponse struct {
	AuthenticationToken string `json:"token"`
	ExpiresAt           string `json:"expires_at"`
}

func NewLoginResponse(authenticationToken string, expiresAt string) *LoginResponse {
	return &LoginResponse{AuthenticationToken: authenticationToken, ExpiresAt: expiresAt}
}

func (req *LoginRequest) setDefaultRegion() {
	if strings.TrimSpace(req.Region) == "" {
		req.Region = config.GetS3Region()
	}
}

func (req *LoginRequest) enforcePath() {
	req.Path = path_utils.EnforceDirectory(req.Path)
}
