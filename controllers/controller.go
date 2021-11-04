package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

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

	// HTMLが崩れてる場合にpanicで終わってしまってもマージさせたいのでdefer。
	defer sitemap.Merge(url.Scheme + "://" + url.Host)
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
	c.String(http.StatusOK, "OK")
}
