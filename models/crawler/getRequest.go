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
	"time"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func GetRequest(r *entity.RequestStruct) {
	fmt.Println("Start GET Request")
	abs := r.Referer.ResolveReference(r.Path)

	if !IsSameOrigin(r, abs) {
		fmt.Println(abs, "is out of Origin.")
		entity.Item.AppendItem(r.Referer.String(), abs.String())
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
			fmt.Fprintln(os.Stderr, "Unable to reach the server.")
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
				} 
				GetRequest(&nextStruct)
			}
    	}
    
		sitemap.Add(*req, (end.Sub(start)).Seconds())
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			fmt.Println("Found: ", abs)
		} else {
			fmt.Println(resp.StatusCode, ": ", abs)
		}
		resp.Body.Close()
		CollectLinks(bytes.NewBuffer(body), abs)
	}	
	fmt.Println(abs, " is Exist.")
}
