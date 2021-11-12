package entity

type Issue struct {
	URL       string
	Kind      string
	Parameter string
	Payload   string
	Request   string
	Response  string
}

var WholeIssue []Issue
