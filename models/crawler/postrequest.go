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
		fmt.Fprintln(os.Stderr, err)
	}

	// Request Header 設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", r.Referer)

	// Requestを投げてResponseを得る。
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to reach the server.")
	}
	
	// fmt.Println("StatusCode =", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// bodyをfunc2に投げる。
	// func2(bytes.NewBuffer(body), base)
	
	return
}