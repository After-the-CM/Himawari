package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Himawari/controllers"
)

func main() {
	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/Himawari", http.Dir("./views/static"))

	router.GET("/api/sitemap", controllers.ReadSitemap)
	router.POST("/api/crawl", controllers.Crawl)
	//	router.POST("/api/deletePath", controller.DeletePath)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/Himawari")
	})
	router.Run(":8080")
}
