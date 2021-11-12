package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"Himawari/models/entity"
	"Himawari/models/logger"
	"Himawari/models/sitemap"
)

func PostRequest(r *entity.RequestStruct) {
	abs := r.Referer.ResolveReference(r.Path)
	if !isSameOrigin(r.Referer, abs) {
		if abs.Scheme == "http" || abs.Scheme == "https" {
			entity.AppendOutOfOrigin(r.Referer.String(), abs.String())
		}
		return
	}

	req, err := http.NewRequest("POST", abs.String(), strings.NewReader(r.Param.Encode()))
	logger.ErrHandle(err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())

	req.PostForm = r.Param

	if !sitemap.IsExist(*req) {
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
					PostRequest(&nextStruct)
				} else {
					GetRequest(&nextStruct)
				}
			}
		}
		CollectLinks(bytes.NewBuffer(body), abs)
	}
}
