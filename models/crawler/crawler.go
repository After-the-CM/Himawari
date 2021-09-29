package crawler

import (
	"log"
	"net/url"

	"Himawari/models/entity"
)

func Crawl(url *url.URL) {
	log.Println("======================================================")
	log.Println("===============     START CRAWLING     ===============")
	log.Println("======================================================")
	log.Println()
	log.Println()
	req := entity.RequestStruct{
		Referer: url,
		Path:    url,
		Param:   url.Query(),
	}
	GetRequest(&req)
}