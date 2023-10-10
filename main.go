package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"monitor/models"
	"monitor/routers"
)

func main() {
	models.InitLogger()
	models.InitConfig()
	models.InitDB()
	models.InitConn()
	models.InitCron().Start()
	gin.SetMode(models.Env.GetString("mode"))
	router := routers.InitRouters()
	server := fmt.Sprintf("%s:%s", models.Env.GetString("server.ip"), models.Env.GetString("server.port"))
	if err := router.Run(server); err != nil {
		models.Logger.Fatalln(err)
	}
}
