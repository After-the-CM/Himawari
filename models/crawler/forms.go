package crawler

import (
	"fmt"
	"net/url"

	"Himawari/models/entity"
)

// memo: Goは配列をconstとして宣言できない。
// patternがあったとしても***Extreme***修正で対応できそう(?)
func PrepareData() (map[string][]string) {
	// https://developer.mozilla.org/ja/docs/Web/HTML/Element/input
	testData := map[string][]string{
		"email": []string{"Himawari@example.com"},
		"url": []string{"http://example.com"},
		"tel": []string{"00012345678", "000-1234-5678", "+81-00-1234-5678"},
		"date": []string{"2020-12-16"},
		"text": []string{"Himawari"},
		"textarea": []string{"Himawari"},
		//"datetime-local"
		
	}
	return testData
}

// func3
func SetValues(form entity.HtmlForm, r entity.RequestStruct) {
	fmt.Println("Start func3")
	testData := PrepareData()

	for i:=0; i<len(form.Values) -2; i++ {
		name := form.Values["Name"][i]
		typ := form.Values["Type"][i]
		tag := form.Values["Tag"][i]
		value := form.Values["Value"][i]

		if value == "NaN" {
			// test dataからtagに対応する値を持ってくる。動かないかもしれない。
			value = testData[typ][0]
		}

		placeholder := form.Values["Placeholder"][0]
		if placeholder != "NaN" {
			value = placeholder
		}
		// requireに対応する処理は未実装
		require := form.Values["Require"][0]
		if require != "NaN" {
			// requireFlag := true
		}
		// patternに対応する処理は未実装
		pattern := form.Values["Pattern"][0]
		if pattern != "NaN" {
			// patternFlag := true 
		}
		
		values := url.Values{}
		values.Set("tag", tag)
		values.Set("type", typ)
		values.Set("name", name)
		values.Set("value", value)
		r.Param = values
		
		// func1へ
		PostRequest(r)
	}
}