package entity

import (
	"net/url"
)

type RequestStruct struct {
	//リンクが存在したページのURL
	Referer *url.URL
	//formの場合はaction
	Path  *url.URL
	Param url.Values
	Form  HtmlForm
}

type HtmlForm struct {
	Action      string
	Method      string
	Type        string
	Name        *string
	Value       *string
	Placeholder *string
	IsOption    bool
	Options     []string
	//Values      url.Values
}

type TestStruct struct {
	//リンクが存在したページのURL
	Origin string
	//formの場合はaction
	Validation string
}
type FoundItemList struct {
	Items map[string][]string
}

var Item = FoundItemList{
	make(map[string][]string),
}

//オリジン外のlinkを収集 存在したページ:link
func (itemList *FoundItemList) AppendItem(place string, u string) {
	//Itemが[string][]stringのため、appendできる。
	itemList.Items[place] = append(itemList.Items[place], u)
}
