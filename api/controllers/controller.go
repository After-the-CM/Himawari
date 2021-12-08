package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	sitemap.CleanSitemapIssue(&entity.JsonNodes)
	c.String(http.StatusOK, "OK")
}

func Crawl(c *gin.Context) {
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "🌻CRAWLを開始します🌻")
	log.Println("===============     START CRAWLING     ===============")
	log.Printf("\n")

	sitemap.Reset()

	var formdata entity.CrawlFormData
	c.Bind(&formdata)

	exclusiveURLs := formdata.ExclusiveURL
	for _, exclusiveURL := range exclusiveURLs {
		u, err := url.Parse(exclusiveURL)
		logger.ErrHandle(err)

		crawler.ExclusiveURLs = append(crawler.ExclusiveURLs, *u)
	}
	delay, err := strconv.Atoi(formdata.Delay)
	if logger.ErrHandle(err) {
		delay = 0
	}
	entity.RequestDelay = time.Duration(delay) * time.Millisecond

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
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "🌻CRAWLが終了しました🌻")
	//sitemap.PrintMap()
	c.String(http.StatusOK, "OK")
}

func ExportOutOfOrigin(c *gin.Context) {
	f := entity.OutOfOrigin
	c.JSON(http.StatusOK, f)
}

var scanflag string = ""

func Scan(c *gin.Context) {
	scanflag = "scanning"
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "🌻SCANを開始します🌻")
	log.Println("===============     START SCANNING     ===============")
	log.Printf("\n")

	var formdata entity.ScanFormData
	c.Bind(&formdata)

	delay, err := strconv.Atoi(formdata.Delay)
	if logger.ErrHandle(err) {
		delay = 0
	}
	entity.RequestDelay = time.Duration(delay) * time.Millisecond

	if formdata.ScanOption == "Quick Scan" {
		scanner.QuickScan = true
	} else {
		scanner.QuickScan = false
	}

	if formdata.LoginURL != "" {
		scanner.SetLoginData(formdata.LoginURL, formdata.LoginReferer, formdata.LoginKey, formdata.LoginValue, formdata.LoginMethod)
	}

	if formdata.LandmarkNumber != 0 {
		scanner.SetGenLandmark(formdata.LandmarkNumber)
	}

	scanner.Scan(&entity.JsonNodes)

	scanflag = "finished"
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "🌻SCANが終了しました🌻")
	c.String(http.StatusOK, "OK")
}

func Scanflag(c *gin.Context) {
	c.String(http.StatusOK, scanflag)
}

func Sort(c *gin.Context) {
	sitemap.SortJson()
	c.String(http.StatusOK, "OK")
}

func Reset(c *gin.Context) {
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "🌻RESETを実行します🌻")
	sitemap.Reset()
	crawler.Reset()
	scanner.Reset()
	c.String(http.StatusOK, "OK")
}

func Readvulns(c *gin.Context) {
	c.JSON(http.StatusOK, entity.Vulnmap)
}
