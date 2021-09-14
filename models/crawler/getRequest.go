package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

//見つけたpathと、refererをくっつけて新しいURLを作る
func GetRequest(r *entity.RequestStruct) {
	fmt.Println("Start GET Request")
	abs := r.Referer.ResolveReference(r.Path)

	//オリジンのチェックを行う
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

	//ヘッダーのセット
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())

	if !sitemap.IsExist(*req) {
		sitemap.Add(*req)

		client := new(http.Client)
		resp, err := client.Do(req)

		if err != nil {
			dump, _ := httputil.DumpRequestOut(req, true)
			fmt.Printf("%s", dump)
			fmt.Fprintln(os.Stderr, "Unable to reach the server.")
		} else {
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				fmt.Println("Found: ", abs)
			} else {
				fmt.Println(resp.StatusCode, ": ", abs)
			}
			//必ずクローズする
			resp.Body.Close()
			CollectLinks(bytes.NewBuffer(body), abs)
		}
	}
	fmt.Println(abs, " is Exist.")
}
