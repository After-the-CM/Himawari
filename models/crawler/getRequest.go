package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"Himawari/models/entity"
	"Himawari/models/sitemap"
)

func GetRequest(r entity.RequestStruct) (forms []entity.HtmlForm) {
	fmt.Println("Start GET Request")
	//Refererは*url.Urlに変更
	//base, _ := url.Parse(r.Referer)
	//Pathは*stringに変更
	rel, _ := url.Parse(*r.Path)
	//abs := base.ResolveReference(rel).String()
	abs := r.Referer.ResolveReference(rel)

	//Pathにabsを入れる必要がないかもしれないと思いコメントアウト化
	//r.Path = abs

	t := entity.TestStruct{
		// Originをhard codingしちゃってる。
		//一度、構造体の型を変更せずに実装してみてる
		Origin:     r.Referer.String(), //"http://localhost:8081/",
		Validation: abs.String(),
	}
	//CheckUrlOrigi→IsSameOrigin(引数も変更)に変更
	if !IsSameOrigin(&r, abs) {
		fmt.Println(abs, "is out of Origin.")
		entity.Item.AppendItem(t.Origin, t.Validation)
		return
	} else {
		fmt.Println(abs)
	}

	/*
		if !CheckUrlOrigin(&t) {
			fmt.Println(abs, "is out of Origin.")
			entity.Item.AppendItem(t.Origin, t.Validation)
			return
		} else {
			fmt.Println(abs)
		}
	*/

	//構造体の変更に伴いString()メソッドの利用に変更
	req, err := http.NewRequest("GET", abs.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		//return
	}
	req.URL.RawQuery = r.Param.Encode()
	req.Header.Set("User-Agent", "Himawari")
	//構造体変更に伴いString()メソッドの利用に変更
	req.Header.Set("Referer", r.Referer.String())

	if !sitemap.IsExist(*req) {
		// fmt.Println("GetRequest:", req)
		sitemap.Add(*req)
		fmt.Println("httpreq-------------------------")
		fmt.Println(*req)

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
			CollectLinks(bytes.NewBuffer(body), abs)
		}
	}
	fmt.Println(abs, " is Exist.")
	return
}
