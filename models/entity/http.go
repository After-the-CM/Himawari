package entity

import (
	"net/http"
	"net/url"
)

type RequestStruct struct {
	Referer string
	Path    string
	Param   url.Values
}

type Message struct {
	Request  *http.Request
	Response *http.Response
}
