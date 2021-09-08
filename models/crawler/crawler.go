package crawler

import (
	"net/url"

	"Himawari/models/entity"
)

func Crawl(url *url.URL) {
	req := entity.RequestStruct{
		Referer: url,         //url.String(),
		Path:    &url.Path,   //url.Path,
		Param:   url.Query(), //url.Query(),
	}
	GetRequest(req)
	//fmt.Println(item.Item)
}
