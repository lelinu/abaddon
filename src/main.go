package main

import "github.com/lelinu/abaddon/src/api/app"

// @title S3 API- GoLelinu
// @version 1.0.0
// @description This is where magic happens
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://golelinu.com/
// @contact.email info@golelinu.com

// @host localhost:8080
// @schemes http

// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main(){
	app.StartApp()
}



