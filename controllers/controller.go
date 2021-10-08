package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"

	"Himawari/models/crawler"
	"Himawari/models/entity"
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
	file, err := c.FormFile("sitemap")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	f, err := file.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	json.Unmarshal(data, &entity.JsonNodes)
	c.String(http.StatusOK, "OK")
}

func Crawl(c *gin.Context) {
	url, _ := url.Parse(c.PostForm("url"))

	// urlのバリデーション

	// HTMLが崩れてる場合にpanicで終わってしまってもマージさせたいのでdefer。
	defer sitemap.Merge(string(url.Scheme) + "://" + string(url.Host))

	crawler.Crawl(url)
	//sitemap.PrintMap()
	c.String(http.StatusOK, "OK")
}

func FoundItem(c *gin.Context) {
	f := entity.Item.Items
	c.JSON(http.StatusOK, f)
}

func Sort(c *gin.Context) {
	sitemap.SortJson()
	c.String(http.StatusOK, "OK")
}
