package authentication

import (
	"github.com/gin-gonic/gin"
	"github.com/lelinu/abaddon/src/api/domain/authentication"
	"github.com/lelinu/abaddon/src/api/services"
	"github.com/lelinu/api_utils/utils/error_utils"
	"net/http"
)

type (
	Controller struct {
		authenticationService services.IAuthenticationService
	}
)

func NewController(AuthenticationService services.IAuthenticationService) *Controller {
	return &Controller{
		authenticationService: AuthenticationService,
	}
}

// Login godoc
// @Summary Endpoint to get jwe token
// @Description  Endpoint to get jwe token
// @Tags Authorization
// @Param request body authentication.LoginRequest true "Get jwe token"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /auth/login [post]
// @Security ApiKeyAuth
func (con *Controller) Login() func(*gin.Context) {

	return func(c *gin.Context) {

		// bind request
		var req authentication.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// login
		resp, err := con.authenticationService.Login(&req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusOK, resp)
		return
	}
}

// RefreshToken godoc
// @Summary Endpoint to refresh a jwe token
// @Description  Endpoint to refresh a jwe token
// @Tags Authorization
// @Param request body authentication.RefreshTokenRequest true "Refresh jwe token"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /auth/refreshtoken [post]
// @Security ApiKeyAuth
func (con *Controller) RefreshToken() func(*gin.Context) {

	return func(c *gin.Context) {

		// bind request
		var req authentication.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// refresh token
		resp, err := con.authenticationService.RefreshToken(&req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusOK, resp)
		return
	}
}


