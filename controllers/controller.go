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

func AddPath(c *gin.Context) {
	path := c.PostForm("path")

	sitemap.AddPath(&entity.Nodes, path)

	c.String(http.StatusOK, "OK")
}
