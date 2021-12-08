package entity

import (
	"net/url"
	"time"
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

var RequestDelay time.Duration = 0

type CrawlFormData struct {
	Name         []string `form:"name[]"`
	Value        []string `form:"value[]"`
	LoginURL     string   `form:"loginURL"`
	LoginReferer string   `form:"loginReferer"`
	LoginKey     []string `form:"loginKey[]"`
	LoginValue   []string `form:"loginValue[]"`
	LoginMethod  []string `form:"loginMethod[]"`
	ExclusiveURL []string `form:"exclusiveURL[]"`
	Delay        string   `form:"delay"`
}

type ScanFormData struct {
	ScanOption     string   `form:"scanOption"`
	LoginURL       string   `form:"loginURL"`
	LoginReferer   string   `form:"loginReferer"`
	LoginKey       []string `form:"loginKey[]"`
	LoginValue     []string `form:"loginValue[]"`
	LoginMethod    []string `form:"loginMethod[]"`
	LandmarkNumber int      `form:"LandmarkNumber"`
	Delay          string   `form:"delay"`
}
