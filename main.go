package main

import (
	"Himawari/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
func main() {
	url, _ := url.Parse("http://localhost:8081/")
	fmt.Println("Start Crawl: ", url)
	crawler.Crawl(url)
}

*/
func main() {
	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/Himawari", http.Dir("./views/static"))

	api := router.Group("/api")
	{
		api.GET("/sitemap", controllers.ReadSitemap)
		api.POST("/crawl", controllers.Crawl)
		api.GET("/found", controllers.FoundItem)
	}
	//router.POST("/api/deletePath", controller.DeletePath)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/Himawari")
	})
	router.Run(":8080")
}
