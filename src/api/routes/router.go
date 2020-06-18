package routes

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/abaddon/src/api/controllers/authentication"
	"github.com/lelinu/abaddon/src/api/controllers/file"
	"github.com/lelinu/abaddon/src/api/controllers/public"
	"github.com/lelinu/abaddon/src/api/middleware"
	"github.com/lelinu/abaddon/src/api/middleware/permissions"
	"github.com/lelinu/abaddon/src/api/providers/awss3_provider"
	"github.com/lelinu/abaddon/src/api/services"
	_ "github.com/lelinu/abaddon/src/docs"
	"github.com/lelinu/api_utils/jwe"
	"github.com/lelinu/api_utils/log/lzap"
	gSwag "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	router = gin.Default()
)

func BuildRouter() *gin.Engine {

	//if api is production switch gin to release mode
	if config.GetApiIsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(
		gin.Recovery(),
		secure.New(secure.Config{
			IsDevelopment:         !config.GetApiIsProd(),
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			ReferrerPolicy:        "strict-origin-when-cross-origin",
		}))

	// use logging service
	logger := lzap.NewService("info", config.GetApiLogPath())
	router.Use(middleware.Logger(logger))

	rg := router.Group(config.GetApiVersionUrl())
	loadSwagger(rg)
	loadControllers(rg, logger)

	return router
}

func loadSwagger(rg *gin.RouterGroup) {
	if config.GetApiIsProd() {
		return
	}

	rg.GET("/swagger/*any", gSwag.WrapHandler(swaggerFiles.Handler))
}

func loadControllers(rg *gin.RouterGroup, logger lzap.IService) {

	// load services
	jweService, err := jwe.NewService(config.GetJweEncryptionAlgorithm(), config.GetJweSecretKey(), config.GetJweIssuer(),
		config.GetJweTokenExpiryInHours(), config.GetJweRefreshTokenExpiryInHours())

	if err != nil{
		panic(err)
	}

	publicController(rg)
	authenticationController(rg, jweService, logger)
	fileController(rg, jweService, logger)
}

func publicController(rg *gin.RouterGroup) {

	var publicController = public.NewController()

	group := rg.Group("public")
	{
		group.GET("/ping", publicController.Ping())
	}
}

func authenticationController(rg *gin.RouterGroup, jweService jwe.IService, logger lzap.IService) {

	// init services
	var provider = awss3_provider.NewAwsS3Provider(config.IsS3LdapEnabled(), config.GetS3EndPoint())
	var authenticationService = services.NewAuthenticationService(logger, provider, jweService)
	var authController = authentication.NewController(authenticationService)

	group := rg.Group("auth")
	{
		group.POST("/login", authController.Login())
		group.POST("/refreshtoken", authController.RefreshToken())
	}
}

func fileController(rg *gin.RouterGroup, jweService jwe.IService, logger lzap.IService) {

	// init services
	var fileService = services.NewFileService(logger)
	var fileController = file.NewController(fileService)
	var authMiddleware = middleware.NewAuthMiddleware(jweService)

	rg.Use(authMiddleware.MiddlewareFunc())
	{
		group := rg.Group("file")
		{
			group.GET("/ls", middleware.AllowClaims(permissions.Owner, permissions.CanRead), fileController.Ls())
			group.GET("/stat", middleware.AllowClaims(permissions.Owner, permissions.CanRead), fileController.Stat())
			group.GET("/cat", middleware.AllowClaims(permissions.Owner, permissions.CanRead), fileController.Cat())
			group.POST("/mkdir", middleware.AllowClaims(permissions.Owner, permissions.CanUpload), fileController.MkDir())
			group.DELETE("/rm", middleware.AllowClaims(permissions.Owner, permissions.CanDelete), fileController.Rm())
			group.POST("/mv", middleware.AllowClaims(permissions.Owner, permissions.CanEdit), fileController.Mv())
			group.POST("/save", middleware.AllowClaims(permissions.Owner, permissions.CanUpload), fileController.Save())
			group.POST("/touch", middleware.AllowClaims(permissions.Owner, permissions.CanUpload), fileController.Touch())
		}
	}
}
