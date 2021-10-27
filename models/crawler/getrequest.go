package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func GetRequest(r *entity.RequestStruct) {
	abs := r.Referer.ResolveReference(r.Path)
	if !isSameOrigin(r.Referer, abs) {
		if abs.Scheme == "http" || abs.Scheme == "https" {
			entity.AppendOutOfOrigin(r.Referer.String(), abs.String())
		}
		return
	}

	req, err := http.NewRequest("GET", abs.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if len(r.Param) != 0 {
		req.URL.RawQuery = r.Param.Encode()
	} else {
		req.URL.RawQuery = abs.RawQuery
	}

	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())

	if !sitemap.IsExist(*req) {
		start := time.Now()
		resp, err := client.Do(req)
		end := time.Now()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		location := resp.Header.Get("Location")
		if location != "" {
			l, _ := url.Parse(location)
			redirect := req.URL.ResolveReference(l)
			if !isSameOrigin(r.Referer, redirect) {
				entity.AppendOutOfOrigin(r.Referer.String(), redirect.String())
				return
			} else {
				nextStruct := entity.RequestStruct{}
				nextStruct.Referer = req.URL
				nextStruct.Path = l
				if resp.StatusCode == 307 {
					nextStruct.Param = r.Param
				}
				GetRequest(&nextStruct)
			}
		}
		sitemap.Add(*req, (end.Sub(start)).Seconds())
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		CollectLinks(bytes.NewBuffer(body), abs)
	}
}
