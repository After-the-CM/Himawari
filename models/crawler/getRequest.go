package crawler

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"net/http/httputil"
	"net/url"
	"bytes"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func GetRequest(r entity.RequestStruct) (forms []entity.HtmlForm) {
	fmt.Println("Start GET Request")
	base, _ := url.Parse(r.Referer)
	rel, _ := url.Parse(r.Path)
	abs := base.ResolveReference(rel).String()
	
	t := entity.TestStruct {
		// Originをhard codingしちゃってる。
		Origin: "http://localhost:8081/",
		Validation: abs,
	}
	if !CheckUrlOrigin(&t) {
		fmt.Println(abs, "is out of Origin.")
		return 
	} else {
		fmt.Println(abs)
	}

	req, err := http.NewRequest("GET", abs, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		//return
	}
	req.URL.RawQuery = r.Param.Encode()
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer)

	if !sitemap.IsExist(*req) {
		// fmt.Println("GetRequest:", req)
		sitemap.Add(*req)

		client := new(http.Client)
		resp, err := client.Do(req)

		if err != nil {
			dump, _ := httputil.DumpRequestOut(req, true)
			fmt.Println(resp.StatusCode)
			fmt.Printf("%s", dump)
			fmt.Fprintln(os.Stderr, "Unable to reach the server.")
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
	fmt.Println(abs, " is Exist.")
	return
}
