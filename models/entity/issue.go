package entity

import "net/url"

type Issue struct {
	URL       string
	Parameter string
	Postparam url.Values
	Getparam  url.Values
	Kind      string
	Request   string
	Response  string
}

var WholeIssue []Issue
