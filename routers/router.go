package routers

import (
	"github.com/gin-gonic/gin"
	"monitor/models"
	"net/http"
)

func InitRouters() (router *gin.Engine) {
	router = gin.Default()
	router.Delims("{[{", "}]}")
	router.LoadHTMLGlob("./static/index.html")
	router.Static("/static", "./static/static")
	router.Static("/img", "./static/img")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"Version": models.GetVersion(),
		})
	})
	router.GET("/list", func(ctx *gin.Context) {
		ok, clients := models.GetData()
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{
				"status": 1,
				"data":   nil,
			})
		}
		ctx.JSONP(http.StatusOK, gin.H{
			"status": 0,
			"data":   clients,
		})
	})

	return router
}
