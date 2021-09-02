package entity

import (
	"net/url"
)

type RequestStruct struct {
	Referer string
	Path    string
	Param   url.Values
}
