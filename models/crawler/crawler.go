package crawler

import (
	"log"
	"net/url"

	"Himawari/models/entity"
)

func Crawl(url *url.URL) {
	log.Println("===============     START CRAWLING     ===============")
	log.Printf("\n")
	req := entity.RequestStruct{
		Referer: url,
		Path:    url,
		Param:   url.Query(),
	}
	GetRequest(&req)
}
