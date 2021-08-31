package crawler

import (
	"fmt"
	"strings"
	"net/http"
	"net/url"

	"Himawari/models/entity"
)

func postRequest(r entity.RequestStruct) {
	// inputのurlはrequestに必要な要素を構造体or引数として渡す。
	// 関数の中で組み合わせてurlを生成する(?)。
	base, _ := url.Parse(r.Referer)
	abs, _ := url.Parse(r.Path)
	s := base.ResolveReference(abs).String()

	// ここでnameの場所に値を入れる。
	postData := r.Param
	
	Request, err := http.NewRequest("POST", s, strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	
	// Request Header 設定
	Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	Request.Header.Set("User-Agent", "Himawari")
	
	// Requestを投げてResponseを得る。
	client := new(http.Client)
	Response, err := client.Do(Request)
	// 関数を抜ける際に必ずResponseをCloseするようにdeferでcloseを呼ぶ。
	defer Response.Body.Close()
	
	if err != nil {
		fmt.Println("Unable to reach the server.")
	} else {
		// bodyをfunc2に投げる。
		// func2(Response)
	}

	// fmt.Println("StatusCode =", Response.StatusCode)
	
	return
}