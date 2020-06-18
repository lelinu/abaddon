package services

import (
	"encoding/json"
	"fmt"
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/abaddon/src/api/domain/authentication"
	"github.com/lelinu/abaddon/src/api/domain/awss3"
	"github.com/lelinu/abaddon/src/api/middleware"
	"github.com/lelinu/abaddon/src/api/providers/awss3_provider"
	"github.com/lelinu/abaddon/src/api/utils/hash_utils"
	"github.com/lelinu/api_utils/jwe"
	"github.com/lelinu/api_utils/log/lzap"
	"github.com/lelinu/api_utils/utils/date_utils"
	"github.com/lelinu/api_utils/utils/error_utils"
)

type AuthenticationService struct {
	logger 		  lzap.IService
	awsS3Provider awss3_provider.IAwsS3Provider
	jweService    jwe.IService
}

type IAuthenticationService interface {
	Login(req *authentication.LoginRequest) (*authentication.LoginResponse, *error_utils.ApiError)
	RefreshToken(req *authentication.RefreshTokenRequest)  (*authentication.RefreshTokenResponse, *error_utils.ApiError)
}

func NewAuthenticationService(logger lzap.IService, awsS3Provider awss3_provider.IAwsS3Provider, jweService jwe.IService) IAuthenticationService {
	return &AuthenticationService{logger: logger, awsS3Provider: awsS3Provider, jweService: jweService}
}

func (s *AuthenticationService) Login(req *authentication.LoginRequest) (*authentication.LoginResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	session := middleware.NewSession(req.Username, req.Password, req.Path, req.EncryptionKey, req.Region,
		date_utils.GetApiCurrentDateTimeString())

	// login
	initRequest := awss3.NewInitRequest(req.Username, req.Password, req.Path, req.Region)
	service, apiErr := s.awsS3Provider.Init(initRequest)
	if apiErr != nil {
		fmt.Printf("api error is %v", apiErr)
		return nil, error_utils.NewUnauthorizedError("Unauthorized")
	}

	lsRequest := &awss3.LsRequest{
		Path: req.Path,
	}
	_, apiErr = service.Ls(lsRequest)
	if apiErr != nil {
		fmt.Printf("api err is %v", apiErr)
		return nil, error_utils.NewUnauthorizedError("Unauthorized")
	}
	// login

	// convert session to bytes
	bytes, err := json.Marshal(session)
	if err != nil {
		return nil, error_utils.NewUnauthorizedError("Unauthorized")
	}

	// encrypt session
	encSession, err := hash_utils.EncryptString(config.GetApiSecretKey(), string(bytes))
	if err != nil {
		return nil, error_utils.NewUnauthorizedError("Unauthorized")
	}

	// generate jwe token
	claims := authentication.CreateNewClaimsForOwner(encSession)
	jweToken, expiry, apiErr := s.jweService.GenerateJweToken(claims.ToMap())
	if apiErr != nil{
		return nil, error_utils.NewUnauthorizedError("Unauthorized")
	}

	// return token and expiry date
	expiryAt := date_utils.ConvertToApiDateFormat(expiry)
	return authentication.NewLoginResponse(jweToken, expiryAt), nil
}

func (s *AuthenticationService) RefreshToken(req *authentication.RefreshTokenRequest)  (*authentication.RefreshTokenResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// call refresh jwe token
	jweToken, expiry, apiErr := s.jweService.RefreshJweToken(req.Token)
	if apiErr != nil{
		return nil, error_utils.NewUnauthorizedError("Unauthorized")
	}

	// return token and expiry date
	expiryAt := date_utils.ConvertToApiDateFormat(expiry)
	return authentication.NewRefreshTokenResponse(jweToken, expiryAt), nil
}
