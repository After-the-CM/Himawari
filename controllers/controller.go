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
		request, _ := http.NewRequest("GET", url.String(), nil)
		entity.Nodes.Messages = append(entity.Nodes.Messages, *request)
		req := entity.RequestStruct{
			Referer: url.String(),
			Path:    url.Path,
			Param:   url.Query(),
		}
		sitemap.PrintMap(entity.Nodes, 0)
		crawler.GetRequest(req)
		sitemap.PrintMap(entity.Nodes, 0)
	*/

	c.String(http.StatusOK, "OK")
}
