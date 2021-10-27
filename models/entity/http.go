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

var OutOfOrigin = make(map[string][]string)

func AppendOutOfOrigin(page string, externalLink string) {
	OutOfOrigin[page] = append(OutOfOrigin[page], externalLink)
}
