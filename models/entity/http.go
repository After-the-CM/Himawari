package entity

import ("net/url")

type RequestStruct struct {
	Referer string
	Path	string
	Param	url.Values
	Form	HtmlForm
}

type HtmlForm struct {
	Action string
	Method string
	Values url.Values
}

type TestStruct struct {
	//リンクが存在したページのURL
	Origin string
	//formの場合はaction
	Validation string
}