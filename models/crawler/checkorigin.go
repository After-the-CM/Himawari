package crawler

import (
	"fmt"
	"net/url"
	"os"

	"Himawari/models/entity"
)

const (
	empty     = ""
	httpPort  = "80"
	httpsPort = "443"
	httpSch   = "http"
	httpsSch  = "https"
)

func IsSameOrigin(r *entity.RequestStruct, n *url.URL) bool {
	switch {
	case (r.Referer.Port() == empty) && (n.Port() == empty):
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
	case r.Referer.Port() == empty:
		if getSchemaPort(&(r.Referer.Scheme), n.Port()) {
			return true
		}
		fallthrough
	case n.Port() == empty:
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
	case httpSch:
		return httpPort == p
	case httpsSch:
		return httpsPort == p
	default:
		fmt.Fprintln(os.Stderr, "http, httpsのスキーム以外のポートは自動解決されません。")
		return false
	}
}

func JudgeMethod(r *entity.RequestStruct) {
	if r.Form.Method == "GET" {
		GetRequest(r)
	} else if r.Form.Method == "POST" {
		PostRequest(r)
	} else {
		fmt.Fprintln(os.Stderr, "Other Methods.")
		return
	}
}
