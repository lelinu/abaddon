package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lelinu/abaddon/src/api/domain/file"
	"github.com/lelinu/abaddon/src/api/middleware"
	"github.com/lelinu/abaddon/src/api/services"
	"github.com/lelinu/api_utils/utils/error_utils"
	"io"
	"net/http"
)

type (
	Controller struct {
		fileService services.IFileService
	}
)

func NewController(fileService services.IFileService) *Controller {
	return &Controller{
		fileService: fileService,
	}
}

// Ls godoc
// @Summary Endpoint to view directories and files
// @Description  Endpoint to view directories and files
// @Tags File
// @Param path query string false "path"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/ls [get]
// @Security ApiKeyAuth
func (con *Controller) Ls() func(*gin.Context) {

	return func(c *gin.Context) {

		// bind request
		var req file.LsRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// ls
		resp, err := con.fileService.Ls(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusOK, resp)
		return
	}
}

// Stat godoc
// @Summary Endpoint to get versions of an object
// @Description  Endpoint to view directories and files
// @Tags File
// @Param path query string false "path"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/stat [get]
// @Security ApiKeyAuth
func (con *Controller) Stat() func(*gin.Context) {

	return func(c *gin.Context) {

		// bind request
		var req file.StatRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// ls
		resp, err := con.fileService.Stat(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusOK, resp)
		return
	}
}

// Cat file godoc
// @Summary Endpoint to download a file
// @Description  Endpoint to download a file
// @Tags File
// @Param request body file.CatRequest true "Download file"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/cat [post]
// @Security ApiKeyAuth
func (con *Controller) Cat() func(*gin.Context) {

	return func(c *gin.Context) {

		// bind request
		var req file.CatRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// ls
		resp, err := con.fileService.Cat(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}
		defer resp.File.Close()

		c.Header("Content-Type", resp.ContentType)
		io.Copy(c.Writer, resp.File)

		return
	}
}

// MkDir godoc
// @Summary Endpoint to create a directory
// @Description  Endpoint to create a directory
// @Tags File
// @Param request body file.MkDirRequest true "Create directory"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/mkdir [post]
// @Security ApiKeyAuth
func (con *Controller) MkDir() func(*gin.Context){

	return func(c *gin.Context) {

		// bind request
		var req file.MkDirRequest
		if err := c.ShouldBind(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// mkdir
		err := con.fileService.MkDir(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusCreated, nil)
		return
	}
}

// Rm file godoc
// @Summary Endpoint to delete a bucket/directory inc files or files
// @Description  Endpoint to delete a bucket/directory inc files or files
// @Tags File
// @Param request body file.RmRequest true "Delete bucket/directory/file"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/rm [delete]
// @Security ApiKeyAuth
func (con *Controller) Rm() func(*gin.Context){

	return func(c *gin.Context) {

		// bind request
		var req file.RmRequest
		if err := c.ShouldBind(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// rm
		err := con.fileService.Rm(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
		return
	}
}

// Mv file godoc
// @Summary Endpoint to move a file
// @Description  Endpoint to move a file
// @Tags File
// @Param request body file.MvRequest true "Move a file"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/mv [post]
// @Security ApiKeyAuth
func (con *Controller) Mv() func(*gin.Context){

	return func(c *gin.Context) {

		// bind request
		var req file.MkDirRequest
		if err := c.ShouldBind(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// mkdir
		err := con.fileService.MkDir(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
		return
	}
}

// Save file godoc
// @Summary Endpoint to save a file
// @Description Endpoint to save a file
// @Tags File
// @Param path formData string true "path"
// @Param file formData file true "file"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/save [post]
// @Security ApiKeyAuth
func (con *Controller) Save() func(*gin.Context){

	return func(c *gin.Context) {

		// bind request
		var req file.SaveRequest
		if err := c.ShouldBind(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)
		fmt.Printf("Scope is %v", scope)

		// save
		resp, err := con.fileService.Save(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusCreated, resp)
		return
	}
}

// Touch file godoc
// @Summary Endpoint to touch a file
// @Description Endpoint to touch a file
// @Tags File
// @Param request body file.TouchRequest true "Create a file without content"
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /file/touch [post]
// @Security ApiKeyAuth
func (con *Controller) Touch() func(*gin.Context){

	return func(c *gin.Context) {

		// bind request
		var req file.TouchRequest
		if err := c.ShouldBind(&req); err != nil {
			restErr := error_utils.NewBadRequestError("bad request")
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		// get scope
		scope := middleware.GetScope(c)

		// save
		err := con.fileService.Touch(scope, &req)
		if err != nil {
			c.JSON(err.HttpStatusCode, err)
			return
		}

		c.JSON(http.StatusCreated, nil)
		return
	}
}





