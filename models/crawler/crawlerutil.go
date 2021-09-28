package crawler

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"Himawari/models/entity"
)


func IsSameOrigin(r *entity.RequestStruct, n *url.URL) bool {
	switch {
	case (r.Referer.Port() == "") && (n.Port() == ""):
		if (r.Referer.Hostname() == n.Hostname()) && (r.Referer.Scheme == n.Scheme) {
			return true
		} else {
			return false
		}
	case r.Referer.Host == n.Host:
		if r.Referer.Scheme == n.Scheme {
			return true
		} else {
			return false
		}
	case r.Referer.Port() == "":
		if getSchemaPort(&(r.Referer.Scheme), n.Port()) {
			return true
		}
		fallthrough
	case n.Port() == "":
		if getSchemaPort(&(n.Scheme), r.Referer.Port()) {
			return true
		}
		return false
	default:
		return false
	}
}

func getSchemaPort(s *string, p string) bool {
	switch *s {
	case "http":
		return (p == "80")
	case "https":
		return (p == "443")
	default:
		fmt.Fprintln(os.Stderr, "http, httpsのスキーム以外のポートは自動解決されません。")
		return false
	}
}

func JudgeMethod(r *entity.RequestStruct) {
	r.Form.Method = strings.ToUpper(r.Form.Method)
	if r.Form.Method == "GET" || r.Form.Method == "" {
		GetRequest(r)
	} else if r.Form.Method == "POST" {
		PostRequest(r)
	} else {
		fmt.Fprintln(os.Stderr, "Other Methods.")
		return
	}
}
