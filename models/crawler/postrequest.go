package crawler

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"net/url"
	"bytes"
	"strings"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func PostRequest(r entity.RequestStruct) {
	base, _ := url.Parse(r.Referer)
	rel, _ := url.Parse(r.Path)
	abs := base.ResolveReference(rel).String()

	postData := r.Param

	req, err := http.NewRequest("POST", abs, strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer)

	if !sitemap.IsExist(*req) {
		// fmt.Println("GetRequest:", req)
		sitemap.Add(*req)

		client := new(http.Client)
		Response, err := client.Do(req)

		if err != nil {
			dump, _ := http.httputil.DumpRequestOut(Request, true)
			fmt.Printf("%s", dump)
			fmt.Println("Unable to reach the server.", err)
		} else {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			func2(bytes.NewBuffer(body), base)
		}
	}
	return
}