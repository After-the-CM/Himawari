package crawler

import (
	"fmt"
	"strings"
	"net/http"
	"net/url"

	"Himawari/models/entity"
)

func PostRequest(r entity.RequestStruct) {
	// inputのurlはrequestに必要な要素を構造体or引数として渡す。
	// 関数の中で組み合わせてurlを生成する(?)。
    base, _ := url.Parse(r.Referer)
    rel, _ := url.Parse(r.Path)
    abs := base.ResolveReference(rel).String()

	// ここでnameの場所に値を入れる。
	postData := r.Param
	
	req, err := http.NewRequest("POST", abs, strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	
	// Request Header 設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	
	// Requestを投げてResponseを得る。
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to reach the server.")
	} 
	resp.Body.Close()
	
	// bodyをfunc2に投げる。
	// func2(Response)

	// fmt.Println("StatusCode =", resp.StatusCode)
	
	return
}