package entity

type Issue struct {
	URL       string
	Kind      string
	Parameter string
	Payload   string
	Cookie    JsonCookie
	Request   string
	Response  string
}

var WholeIssue []Issue
