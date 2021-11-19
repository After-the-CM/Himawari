package crawler

import (
	"net/url"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func Crawl(url *url.URL) {
	// HTMLが崩れてる場合にpanicで終わってしまってもマージさせたいのでdefer。
	defer sitemap.Merge(url, jar)

	req := entity.RequestStruct{
		Referer: url,
		Path:    url,
		Param:   url.Query(),
	}
	GetRequest(&req)
}

func Reset() {
	applyData = map[string]string{}
}
