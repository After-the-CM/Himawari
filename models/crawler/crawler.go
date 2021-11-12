package crawler

import (
	"net/url"

	"Himawari/models/entity"
)

func Crawl(url *url.URL) {
	req := entity.RequestStruct{
		Referer: url,
		Path:    url,
		Param:   url.Query(),
	}
	GetRequest(&req)
}
