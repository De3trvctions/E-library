package main

import (
	"e-library/config"
	initilize "e-library/initialize"
	_ "e-library/routers"
	"e-library/validation"
	"fmt"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	initilize.InitLogs()
	initilize.InitDB()
	validation.Init()

	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Token", "Language"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}), web.WithReturnOnOutput(true))

	// To enable getting json body for the request instead of only using form-data
	web.BConfig.CopyRequestBody = true

	web.Run(fmt.Sprintf(":%d", config.HttpPort))
}
