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
	fmt.Println("Start POST Request")
	abs := r.Referer.ResolveReference(r.Path)

	if !IsSameOrigin(r, abs) {
		fmt.Println(abs, "is out of Origin.")
		entity.Item.AppendItem(r.Referer.String(), abs.String())
		return
	} else {
		fmt.Println(abs)
	}
	req, err := http.NewRequest("POST", abs.String(), strings.NewReader(r.Param.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	//ヘッダーのセット
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())

	req.PostForm = r.Param

	if !sitemap.IsExist(*req) {

		start := time.Now()
    client := &http.Client{
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
				fmt.Println(redirect, "is out of Origin.")
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
		if resp.StatusCode == 200 {
			fmt.Println("Found: ", abs)
		} else {
			fmt.Println(resp.StatusCode, ": ", abs)
		}
		//必ずクローズする
		resp.Body.Close()
		//次のlinkを探す
		CollectLinks(bytes.NewBuffer(body), abs)
	}
}
