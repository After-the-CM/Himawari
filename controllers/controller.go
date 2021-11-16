package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"Himawari/models/crawler"
	"Himawari/models/entity"
	"Himawari/models/logger"
	"Himawari/models/scanner"
	"Himawari/models/sitemap"
)

func ReadSitemap(c *gin.Context) {
	c.JSON(http.StatusOK, entity.JsonNodes)
}

func DownloadSitemap(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=sitemap.json")
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.JSON(http.StatusOK, entity.JsonNodes)
}

func DownloadMarkdown(c *gin.Context) {
	t := time.Now()
	c.Header("Content-Disposition", "attachment; filename=Himawari_Report_"+t.Format("2006-01-02")+".md")
	c.Header("Content-Type", "text/markdown; charset=UTF-8")
	md := scanner.MarkDown()
	c.String(http.StatusOK, md)
}

func UploadSitemap(c *gin.Context) {
	sitemap.Reset()

	file, err := c.FormFile("sitemap")
	logger.ErrHandle(err)
	f, err := file.Open()
	logger.ErrHandle(err)
	data, err := io.ReadAll(f)
	logger.ErrHandle(err)
	json.Unmarshal(data, &entity.JsonNodes)
	c.String(http.StatusOK, "OK")
}

func Crawl(c *gin.Context) {
	log.Println("===============     START CRAWLING     ===============")
	log.Printf("\n")

	sitemap.Reset()

	url, err := url.Parse(c.PostForm("url"))
	if logger.ErrHandle(err) {
		return
	}

	// urlのバリデーション

	crawler.Crawl(url)
	//sitemap.PrintMap()
	c.String(http.StatusOK, "OK")
}

func ExportOutOfOrigin(c *gin.Context) {
	f := entity.OutOfOrigin
	c.JSON(http.StatusOK, f)
}

func Scan(c *gin.Context) {
	log.Println("===============     START SCANNING     ===============")
	log.Printf("\n")
	scanner.Scan(&entity.JsonNodes)
	c.JSON(http.StatusOK, entity.WholeIssue)
}

func Sort(c *gin.Context) {
	sitemap.SortJson()
	c.String(http.StatusOK, "OK")
}

func Reset(c *gin.Context) {
	sitemap.Reset()
	scanner.Reset()
	c.String(http.StatusOK, "OK")
}

func Readvulns(c *gin.Context) {
	c.JSON(http.StatusOK, entity.Vulnmap)
}
