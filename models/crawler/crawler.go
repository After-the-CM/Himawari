package crawler

import (
	"net/url"
	
	"Himawari/models/entity"
)

func Crawl(url *url.URL) {
	req := entity.RequestStruct{
		Referer: url.String(),
		Path: url.Path,
		Param: url.Query(),
	}
	GetRequest(req)
}