package crawler

import (
	"bytes"
	"fmt"
	"log"
	"io"
	"net/url"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func PostRequest(r *entity.RequestStruct) {
	abs := r.Referer.ResolveReference(r.Path)
	if !IsSameOrigin(r, abs) {
		if abs.Scheme == "http" || abs.Scheme == "https" {
			entity.Item.AppendItem(r.Referer.String(), abs.String())
			return
		}
	}

	req, err := http.NewRequest("POST", abs.String(), strings.NewReader(r.Param.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())

	req.PostForm = r.Param

	if !sitemap.IsExist(*req) {
		start := time.Now()
		client := &http.Client {
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		log.Println(req)
		resp, err := client.Do(req)
		log.Println(resp)
		end := time.Now()

		if err != nil {
			dump, _ := httputil.DumpRequestOut(req, true)
			fmt.Printf("%s", dump)
			fmt.Println("Unable to reach the server.", err)
		}

		location := resp.Header.Get("Location")
		if location != "" {
			l, _ := url.Parse(location)
			redirect := r.Referer.ResolveReference(l)
			if !IsSameOrigin(r, redirect) {
				entity.Item.AppendItem(r.Referer.String(), redirect.String())
				return
			} else {
				nextStruct := entity.RequestStruct{}
				nextStruct.Referer = r.Referer
				nextStruct.Path = l
				if resp.StatusCode == 307 {
					nextStruct.Param = r.Param
					PostRequest(&nextStruct)
				} else {
					GetRequest(&nextStruct)
				} 
			}
		}
		sitemap.Add(*req, (end.Sub(start)).Seconds())
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		CollectLinks(bytes.NewBuffer(body), abs)
	}
}
