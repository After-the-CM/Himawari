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

var loginMsg = entity.JsonMessage{
	URL:       "http://localhost:18080/osci/login.php",
	Referer:   "http://localhost:18080/osci/login.php",
	GetParams: url.Values{},
	PostParams: url.Values{
		"name": []string{"yoden"},
		"pass": []string{"pass"},
	},
}

func login(jar http.CookieJar) http.CookieJar {
	var client4login = &http.Client{
		Jar: jar,
		/*
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		*/
		Transport: logger.LoggingRoundTripper{
			Proxied: http.DefaultTransport,
		},
	}

	var req *http.Request
	if len(loginMsg.PostParams) != 0 {
		req = genPostParamReq(&loginMsg, &loginMsg.PostParams)
	} else {
		req = genGetParamReq(&loginMsg, &loginMsg.GetParams)
	}

	/*_, err :=*/
	client.Do(req)
	// logger.ErrHandle(err)
	return client4login.Jar
}

func createGetReq(url string, ref string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Himawari")
	if ref != "" {
		req.Header.Set("Referer", ref)
	}
	return req
}

func createPostReq(url string, ref string, p url.Values) *http.Request {
	req, _ := http.NewRequest("POST", url, strings.NewReader(p.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", ref)
	return req
}

func genGetParamReq(j *entity.JsonMessage, gp *url.Values) *http.Request {
	req := createGetReq(j.URL, j.Referer)
	req.URL.RawQuery = gp.Encode()
	return req
}

func genPostParamReq(j *entity.JsonMessage, pp *url.Values) *http.Request {
	req := createPostReq(j.URL, j.Referer, *pp)
	req.URL.RawQuery = j.GetParams.Encode()
	return req
}
