package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func ReadSitemap(c *gin.Context) {
	j := sitemap.MtoJ(entity.Nodes)
	c.JSON(http.StatusOK, j)
}

func Crawl(c *gin.Context) {
	/*
		url, _ := url.Parse(c.PostForm("url"))

		// urlのバリデーション

		 crawler.Crawl(url)
	*/
	sitemap.PrintMap()
	c.String(http.StatusOK, "OK")
}
