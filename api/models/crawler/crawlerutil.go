package crawler

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

var jar, _ = cookiejar.New(nil)
var client = &http.Client{
	Jar: jar,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Transport: logger.LoggingRoundTripper{
		Proxied: http.DefaultTransport,
	},
}

func isSameOrigin(ref *url.URL, loc *url.URL) bool {
	rport, lport := ref.Port(), loc.Port()
	if rport == "" {
		rport = getSchemaPort(ref.Scheme)
	}
	if lport == "" {
		lport = getSchemaPort(loc.Scheme)
	}
	if ref.Hostname() == loc.Hostname() && rport == lport && ref.Scheme == loc.Scheme {
		return true
	}
	return false
}

func getSchemaPort(s string) string {
	switch s {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		fmt.Fprintln(os.Stderr, "http, httpsのスキーム以外のポートは自動解決されません。")
		return ""
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
