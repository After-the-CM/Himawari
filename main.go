package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mbsd/controllers"
)

func main() {
	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/mbsd", http.Dir("./views/static"))

	router.GET("/api/sitemap", controllers.ReadSitemap)
	router.POST("/api/addPath", controllers.AddPath)
	//	router.POST("/api/deletePath", controller.DeletePath)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/mbsd")
	})
	router.Run(":8080")
}
