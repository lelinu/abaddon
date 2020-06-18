package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lelinu/api_utils/log/lzap"
	"time"
)

func Logger(logger lzap.IService) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now().UTC()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error("Logger middleware", errors.New(e), "")
			}
		} else {

			logger.Info(path, fmt.Sprintf("status:%v", c.Writer.Status()),
				fmt.Sprintf("host:%v", c.Request.Host),
				fmt.Sprintf("method:%v", c.Request.Method),
				fmt.Sprintf("path:%v", path),
				fmt.Sprintf("query:%v", query),
				fmt.Sprintf("ip:%v", c.ClientIP()),
				fmt.Sprintf("referer:%v", c.Request.Referer()),
				fmt.Sprintf("user-agent:%v", c.Request.UserAgent()),
				fmt.Sprintf("latency:%v", latency),
			)
		}
	}
}
