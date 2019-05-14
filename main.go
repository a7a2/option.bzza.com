package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"option.bzza.com/controllers"
	"option.bzza.com/system"
)

func main() {
	//currentDirectory := system.GetCurrentDirectory()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	//router.LoadHTMLGlob("templates/default")
	router.Static("/templates", "./templates")
	router.LoadHTMLFiles("./templates/default/index.html")

	router.Static("/css", "./templates/default/css")
	router.Static("/res/img", "./templates/default/res/img")
	router.Static("/scripts", "./templates/default/js")
	router.Static("/settings", "./templates/default/settings")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.POST("/ver", controllers.PostVer)
	router.POST("/reg", controllers.PostReg)
	router.GET("/auth/:loginNum/:password", controllers.GetLogin)

	router.GET("/session/:token", func(c *gin.Context) {
		controllers.Wshandler(c.Writer, c.Request)
	})

	router.POST("/query", controllers.Query)
	if err := router.RunTLS(":443", system.Conf.Certificate.Pem, system.Conf.Certificate.Key); err != nil {
		log.Println(err)
		router.Run(":" + system.Conf.Web.Port)
	}
}
