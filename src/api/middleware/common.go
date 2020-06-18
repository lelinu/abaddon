package middleware

import "github.com/gin-gonic/gin"

//extractClaims will extract claims
func extractClaims(c *gin.Context) map[string]interface{} {
	claims, exists := c.Get("JWE_PAYLOAD")
	if !exists {
		return map[string]interface{}{}
	}

	return claims.(map[string]interface{})
}
