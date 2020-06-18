package app

import (
	"fmt"
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/abaddon/src/api/routes"
)

func StartApp(){

	router := routes.BuildRouter()
	router.Run(fmt.Sprintf(":%v", config.GetApiPort()))
}