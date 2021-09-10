package entity

import (
	"net/url"
)

/*
type RequestStruct struct {
	Referer string
	Path    string
	Param   url.Values
	Form    HtmlForm
}
*/

type RequestStruct struct {
	//リンクが存在したページのURL
	Referer *url.URL
	//formの場合はaction
	Path *url.URL
	//Path *string

	Param url.Values

	//Method string
	//Form []HtmlForm
	Form HtmlForm
}

type HtmlForm struct {
	Action      string
	Method      string
	Type        string
	Name        string
	Value       string
	Placeholder string
	IsOption    bool
	//Values      url.Values
}
type TestStruct struct {
	//リンクが存在したページのURL
	Origin string
	//formの場合はaction
	Validation string
}
type FoundItemList struct {
	//Place string
	//Item string
	Items map[string][]string
}

var Item = FoundItemList{
	make(map[string][]string),
}

func (itemList *FoundItemList) AppendItem(place string, u string) {
	//Itemが[string][]stringのため、appendできる。
	itemList.Items[place] = append(itemList.Items[place], u)
}
