package main

import (
	"Himawari/controllers"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

/*
func main() {
	proxyUrl, err := url.Parse("http://172.16.82.190:8001")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	url, _ := url.Parse("http://localhost:8081/")
	fmt.Println("Start Crawl: ", url)
	crawler.Crawl(url)

}
*/

func main() {
	proxyUrl, err := url.Parse("http://172.16.82.190:8001")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/Himawari", http.Dir("./views/static"))

	api := router.Group("/api")
	{
		api.GET("/sitemap", controllers.ReadSitemap)
		api.POST("/crawl", controllers.Crawl)
		api.GET("/found", controllers.FoundItem)
	}

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/Himawari")
	})
	router.Run(":8080")
}
