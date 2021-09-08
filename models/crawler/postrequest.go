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
	//GetRequestと同じように変更
	//base, _ := url.Parse(r.Referer)
	rel, _ := url.Parse(*r.Path)
	//abs := base.ResolveReference(rel).String()
	//abs := r.Referer.ResolveReference(rel).String()
	abs := r.Referer.ResolveReference(rel)

	//構造体の変更に伴いString()メソッドの利用に変更
	t := entity.TestStruct{
		// Originをhard codingしちゃってる。
		Origin:     r.Referer.String(), //"http://localhost:8081/",
		Validation: abs.String(),
	}
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
			return
		} else {
			fmt.Println(abs)
		}
	*/

	postData := r.Param

	//構造体の変更に伴いString()メソッドの利用に変更
	req, err := http.NewRequest("POST", abs.String(), strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer.String())
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
			//Refererではなく、新たなアクセス先だと思うのでabsに変更
			//CollectLinks(bytes.NewBuffer(body), base)
			CollectLinks(bytes.NewBuffer(body), abs)
		}
	}
	//return
}
