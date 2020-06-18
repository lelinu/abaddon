package middleware

import (
	"errors"
	"github.com/lelinu/api_utils/jwe"
	"strings"

	"github.com/gin-gonic/gin"

	"net/http"
)

type AuthMiddleware struct {
	jweService            jwe.IService
	unauthorizedFunc      func(*gin.Context, int, string)
	identityHandlerFunc   func(*gin.Context) interface{}
	httpStatusMessageFunc func(e error, c *gin.Context) string
	authorisationFunc     func(data interface{}, c *gin.Context) bool
}

func NewAuthMiddleware(jweService jwe.IService) *AuthMiddleware {
	return &AuthMiddleware{
		jweService: jweService,
	}
}

var (
	tokenLookup   = "header:Authorization"
	tokenHeadName = "Bearer"
	identityKey   = "id"

	errForbidden         = errors.New("you don't have permission to access this resource")
	errEmptyAuthHeader   = errors.New("auth header is empty")
	errInvalidAuthHeader = errors.New("auth header is invalid")
)

//MiddlewareFunc function
func (a *AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		a.validate(c)
	}
}

//validate will validate the token
func (a *AuthMiddleware) validate(c *gin.Context) {

	//init middleware
	a.init()

	//parse token
	token, err := a.parseToken(c)
	if err != nil {
		a.unauthorized(c, http.StatusForbidden, a.httpStatusMessageFunc(errForbidden, c))
		return
	}

	//validate jwe token
	claims, apiErr := a.jweService.ValidateJweToken(token)
	if apiErr != nil {
		a.unauthorized(c, http.StatusForbidden, a.httpStatusMessageFunc(errForbidden, c))
		return
	}

	c.Set("JWE_PAYLOAD", claims)
	identity := a.identityHandlerFunc(c)

	if identity != nil {
		c.Set(identityKey, identity)
	}

	c.Next()
}

//init will init the middleware functionality
func (a *AuthMiddleware) init() error {

	if a.unauthorizedFunc == nil {
		a.unauthorizedFunc = func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		}
	}

	if a.identityHandlerFunc == nil {
		a.identityHandlerFunc = func(c *gin.Context) interface{} {
			claims := extractClaims(c)
			return claims[identityKey]
		}
	}

	if a.httpStatusMessageFunc == nil {
		a.httpStatusMessageFunc = func(e error, c *gin.Context) string {
			return e.Error()
		}
	}

	return nil
}

//unauthorized used to return unauthorized from gin
func (a *AuthMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.Abort()
	a.unauthorizedFunc(c, code, message)
}

//parseToken will parse the token
func (a *AuthMiddleware) parseToken(c *gin.Context) (string, error) {
	var token string
	var err error

	methods := strings.Split(tokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = a.tokenFromHeader(c, v)
		}
	}

	if err != nil {
		return "", err
	}

	return token, nil
}

//tokenFromHeader will read the jwt from header
func (a *AuthMiddleware) tokenFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == tokenHeadName) {
		return "", errInvalidAuthHeader
	}

	return parts[1], nil
}
