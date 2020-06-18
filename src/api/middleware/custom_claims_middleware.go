package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lelinu/api_utils/utils/error_utils"
)

// AllowClaims
func AllowClaims(claims ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		validateClaims(c, claims...)
		return
	}
}

func validateClaims(c *gin.Context, claims ...string) {

	// token claims
	tokenContent := extractClaims(c)
	if value, ok := tokenContent["is_owner"].(bool); ok {
		if value == true{
			c.Next()
			return
		}
	}

	// get permissions
	permissions, ok := tokenContent["permissions"].(map[string]interface{})
	if !ok{
		unauthorized(c)
		return
	}

	// validating claims against token claims
	for _, claim := range claims {
		for k, v := range permissions {
			if k == claim && v == true {
				c.Next()
				return
			}
		}
	}

	// no claim match
	unauthorized(c)
	return
}

//unauthorized will return unauthorized to the client
func unauthorized(c *gin.Context){

	response := error_utils.NewUnauthorizedError("unauthorized")
	c.Writer.WriteHeader(response.HttpStatusCode)
	c.Writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(response)
	c.Abort()
}



