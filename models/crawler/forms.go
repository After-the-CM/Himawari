package crawler

import (
	"fmt"
	"net/url"

	"Himawari/models/entity"
)

// patternがあったとしても***Extreme***修正で対応できそう(?)
//func PrepareData() map[string]string {
// https://developer.mozilla.org/ja/docs/Web/HTML/Element/input
// 実装が重くなったので、テストケースは1種類で実装。
var TestData = map[string]string{
	"email":    "Himawari@example.com",
	"url":      "http://example.com",
	"tel":      "00012345678",
	"date":     "2020-12-16",
	"text":     "Himawari",
	"textarea": "Himawari",
	//"datetime-local"
}

//return testData
//}

// func3
func SetValues(form []entity.HtmlForm, r entity.RequestStruct) {
	fmt.Println("Start func3")
	//testData := PrepareData()

	r.Form.Action = form[0].Action
	path, _ := url.Parse(form[0].Action)
	r.Path = path
	r.Form.Method = form[0].Method
	fmt.Println("form name & form method ::", r.Form.Method, form[0].Method)
	//attrs := make(map[int](map[string]string), len(form.Values["Name"]))

	/*
		for i := 0; i < len(form.Values["Name"]); i++ {
			attr := make(map[string]string, len(form.Values))
			for j, v := range form.Values {
				attr[j] = v[i]
			}
			attrs[i] = attr
		}
	*/

	//var inputs []url.Values
	fmt.Println("!!!", form)
	for _, v := range form {
		values := url.Values{}
		switch {
		case v.Type != "submit":
			if v.Value == "" {
				values.Add(v.Name, TestData[v.Name])
				fmt.Println(v.Name)
				fmt.Println(TestData[v.Name])
			} else {
				values.Add(v.Name, v.Value)
				fmt.Println("v.name&v.value", v.Name, v.Value, values)
			}
			fmt.Println("NAME     ", v.Name, "VALUEEE   ", values)

		}
		//inputs = append(inputs, values)
		if len(values) != 0 {
			r.Param = values
		}
		fmt.Println("r.parammmmmm", r.Param)
	}
	fmt.Println(r)
	if r.Form.Method == "GET" || r.Form.Method == "get" {
		GetRequest(r)
	} else if r.Form.Method == "POST" || r.Form.Method == "post" {
		fmt.Println("りくえすとおおおお", r)
		PostRequest(r)
	} else {
		return
	}
	/*
		for i := 0; i < len(form.Values["Name"]); i++ {
			values := url.Values{}
			for j, v := range attrs[i] {
				switch j {
				case "Tag":
					if v != "NaN" {
						values.Set("tag", attrs[i][j])
					}
				case "Type":
					if v != "NaN" {
						values.Set("type", attrs[i][j])
						if attrs[i]["Value"] == "NaN" {
							values.Set("value", TestData[attrs[i][j]])
						}
					}
				case "Name":
					if v != "NaN" {
						values.Set("name", attrs[i][j])
					}
				case "Value":
					if v == "NaN" {
						values.Set("value", TestData[attrs[i]["Type"]])
					} else {
						values.Set("value", attrs[i][j])
					}
				case "Placeholder":
					if v != "NaN" {
						attrs[i]["Value"] = v
						values.Set("value", attrs[i]["Values"])
					}
				// patternに対応する処理は未実装
				case "Pattern":
					fmt.Println("DETECTED!! HTML attribute: Pattern!!")
					// requireに対応する処理は未実装だが、ValueにはNaN以外の値を入れて送信してるからいいかな
				}
			}
			r.Param = values
			PostRequest(r)
		}
	*/
}
