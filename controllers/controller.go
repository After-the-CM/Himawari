package controllers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"Himawari/models/crawler"
	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func ReadSitemap(c *gin.Context) {
	j := sitemap.MtoJ(entity.Nodes)
	c.JSON(http.StatusOK, j)
}

func Crawl(c *gin.Context) {

	url, _ := url.Parse(c.PostForm("url"))

	// urlのバリデーション

	crawler.Crawl(url)

	sitemap.PrintMap()
	c.String(http.StatusOK, "OK")
}

func FoundItem(c *gin.Context) {
	f := entity.Item.Items
	c.JSON(http.StatusOK, f)
}
