package crawler

import (
	"fmt"
	"net/url"

	"Himawari/models/entity"
)

// memo: Goは配列をconstとして宣言できない。
// patternがあったとしても***Extreme***修正で対応できそう(?)
// 実装が重くなったので、テストケースは1種類で実装。
var TestData = map[string]string{
	"email":    "Himawari@example.com",
	"url":      "http://example.com",
	"tel":      "00012345678",
	"date":     "2020-12-16",
	"text":     "Himawari",
	"textarea": "Himawari",
	"input":    "I am Himawari",
}

// func3
func SetValues(form []entity.HtmlForm, r *entity.RequestStruct) {
	fmt.Println("Start func3")

	r.Form.Action = form[0].Action
	path, _ := url.Parse(form[0].Action)
	r.Path = path
	r.Form.Method = form[0].Method

	values := url.Values{}
	for _, v := range form {
		if v.Name != nil {
			switch {
			//selectの場合(ほかのものより先に実行しないとうまくいかなかった)
			case v.IsOption:
				if len(v.Options) == 0 {
					values.Set(*v.Name, v.Options[0])
				} else {
					values.Set(*v.Name, v.Options[1])
				}

			//submitではないもの
			case v.Type != "submit":
				//placeholderがあるのならば
				if v.Placeholder != nil {
					values.Set(*v.Name, *v.Placeholder)
				} else if v.Value == nil { //valueが空っぽな場合はテストデータを取得
					values.Set(*v.Name, TestData[*v.Name])
				} else {
					values.Set(*v.Name, *v.Value)
				}
			}

		}

		if len(values) != 0 {
			r.Param = values
		}
	}
	//リクエストの送信
	JudgeMethod(r)
}
