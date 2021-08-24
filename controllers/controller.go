package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

var node = entity.Node{
	Path: "/",
}

func ReadSitemap(c *gin.Context) {
	j := sitemap.MtoJ(node)
	c.JSON(http.StatusOK, j)
}

func AddPath(c *gin.Context) {
	path := c.PostForm("path")

	sitemap.AddPath(&node, path)

	c.String(http.StatusOK, "OK")

}
