package controllers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"Himawari/models/crawler"
	"Himawari/models/sitemap"
)

func ReadSitemap(c *gin.Context) {
	j := sitemap.Json()
	c.JSON(http.StatusOK, j)
}

func Crawl(c *gin.Context) {
	url, _ := url.Parse(c.PostForm("url"))

	// urlのバリデーション

	crawler.Crawl(url)
	sitemap.PrintMap()
	c.String(http.StatusOK, "OK")
}
