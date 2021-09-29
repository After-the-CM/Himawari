package main

import (
	"fmt"
	"io"
	"os"

	"log"
	"net/http"
	"net/url"
	"time"

	"Himawari/controllers"

	"github.com/gin-gonic/gin"
)

func loggingSetting() {
	layout := "2006-01-02_15:04:05"
	dirName := "log"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.Mkdir(dirName, 0666)
	}
	t := time.Now()
	fileName := "log/" + t.Format(layout) + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetFlags(log.Flags() &^ log.LstdFlags)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func init() {
	loggingSetting()
}

func main() {
	/*
		proxyUrl, err := url.Parse("http://172.16.82.190:8001")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	*/

	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/Himawari", http.Dir("./views/static"))

	api := router.Group("/api")
	{
		api.GET("/sitemap", controllers.ReadSitemap)
		api.POST("/crawl", controllers.Crawl)
		api.GET("/found", controllers.FoundItem)
		api.GET("/sort", controllers.Sort)
	}
	//	router.POST("/api/deletePath", controller.DeletePath)
	router.GET("/download", controllers.DownloadSitemap)
	router.POST("/upload", controllers.UploadSitemap)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/Himawari")
	})
	router.Run(":8080")
}
