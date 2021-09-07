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

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func PostRequest(r entity.RequestStruct) {
	fmt.Println("Start POST Request")
	base, _ := url.Parse(r.Referer)
	rel, _ := url.Parse(r.Path)
	abs := base.ResolveReference(rel).String()

	t := entity.TestStruct{
		// Originをhard codingしちゃってる。
		Origin:     "http://localhost:8081/",
		Validation: abs,
	}
	if !CheckUrlOrigin(&t) {
		fmt.Println(abs, "is out of Origin.")
		return
	} else {
		fmt.Println(abs)
	}

	postData := r.Param

	req, err := http.NewRequest("POST", abs, strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer)
	req.PostForm = r.Param
	if !sitemap.IsExist(*req) {
		// fmt.Println("GetRequest:", req)
		sitemap.Add(*req)

		client := new(http.Client)
		resp, err := client.Do(req)

		if err != nil {
			dump, _ := httputil.DumpRequestOut(req, true)
			fmt.Printf("%s", dump)
			fmt.Println("Unable to reach the server.", err)
		} else {
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				fmt.Println("Found: ", abs)
			} else {
				fmt.Println(resp.StatusCode, ": ", abs)
			}
			resp.Body.Close()
			CollectLinks(bytes.NewBuffer(body), base)
		}
	}
	return
}
