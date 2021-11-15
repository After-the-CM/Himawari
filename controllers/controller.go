package controllers

import (
	"encoding/json"
	"fmt"
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

	var formdata entity.CrawlFormData
	c.Bind(&formdata)

	crawler.SetApplydata(formdata.Name, formdata.Value)
	if formdata.LoginURL != "" {
		crawler.SetLoginData(formdata.LoginURL, formdata.LoginReferer, formdata.LoginKey, formdata.LoginValue, formdata.LoginMethod)
	}

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
	scanflag = "scanning"
	log.Println("===============     START SCANNING     ===============")
	log.Printf("\n")

	var formdata entity.ScanFormData
	c.Bind(&formdata)

	fmt.Println("scanner.QuickScan is ", scanner.QuickScan)
	fmt.Println(formdata)
	if formdata.ScanOption == "Quick Scan" {
		scanner.QuickScan = true
		fmt.Println("scanner.QuickScan is ", scanner.QuickScan)
	}

	if formdata.LoginURL != "" {
		scanner.SetLoginData(formdata.LoginURL, formdata.LoginReferer, formdata.LoginKey, formdata.LoginValue, formdata.LoginMethod)
	}

	fmt.Println("form.Randmarrrrrrk", formdata.RandmarkNumber)
	if formdata.RandmarkNumber != 0 {
		scanner.SetGenRandmark(formdata.RandmarkNumber)
		fmt.Println("form.Randmarrrrrrk", formdata.RandmarkNumber)
	}

	scanner.Scan(&entity.JsonNodes)
	//c.JSON(http.StatusOK, entity.WholeIssue)
	scanflag = "finished"
	c.String(http.StatusOK, "OK")
}

var scanflag string = ""

func Scanflag(c *gin.Context) {
	c.String(http.StatusOK, scanflag)
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
