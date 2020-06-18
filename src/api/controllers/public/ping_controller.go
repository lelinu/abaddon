package public

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	Controller struct {

	}
)

func NewController() *Controller {
	return &Controller{
	}
}

func (con *Controller) Ping() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
		return
	}
}