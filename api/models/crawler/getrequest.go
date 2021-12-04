package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"Himawari/models/entity"
	"Himawari/models/logger"
	"Himawari/models/sitemap"
)

func GetRequest(r *entity.RequestStruct) {
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	abs := r.Referer.ResolveReference(r.Path)
	fmt.Println("GET", abs)
	if !isSameOrigin(r.Referer, abs) {
		if abs.Scheme == "http" || abs.Scheme == "https" {
			entity.AppendOutOfOrigin(r.Referer.String(), abs.String())
		}
		return
	}

	req, err := http.NewRequest("GET", abs.String(), nil)
	logger.ErrHandle(err)

	if len(r.Param) != 0 {
		req.URL.RawQuery = r.Param.Encode()
	} else {
		req.URL.RawQuery = abs.RawQuery
	}

	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())

	if !sitemap.IsExist(*req) {
		time.Sleep(entity.RequestDelay)

		start := time.Now()
		resp, err := client.Do(req)
		end := time.Now()

		if logger.ErrHandle(err) {
			dump, err := httputil.DumpRequestOut(req, true)
			logger.ErrHandle(err)
			fmt.Fprintln(os.Stderr, string(dump))
			return
		}

		sitemap.Add(*req, (end.Sub(start)).Seconds())

		body, err := io.ReadAll(resp.Body)
		logger.ErrHandle(err)

		defer resp.Body.Close()

		location := resp.Header.Get("Location")
		if location != "" {
			//locationのParseができないとリダイレクトができないためreturn
			l, err := url.Parse(location)
			if logger.ErrHandle(err) {
				return
			}
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
		CollectLinks(bytes.NewBuffer(body), abs)
	}
}
