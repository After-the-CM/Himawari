package entity

import "net/url"

type Issue struct {
	URL       string
	Kind      string
	Parameter string
	Cookie    JsonCookie
	Getparam  url.Values
	Postparam url.Values
	Request   string
	Response  string
}

var WholeIssue []Issue
