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

func SetApplydata(name []string, value []string) {
	for i := 0; len(name) > i; i++ {
		applyData[name[i]] = value[i]
	}
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
	GetParams:  url.Values{},
	PostParams: url.Values{},
}

func SetLoginData(url string, ref string, keys []string, values []string, methods []string) {

	loginMsg.URL = url
	loginMsg.Referer = ref

	for i := 0; len(keys) > i; i++ {
		if methods[i] == "GET" {
			loginMsg.GetParams.Set(keys[i], values[i])
		} else {
			loginMsg.PostParams.Set(keys[i], values[i])
		}
	}

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
	var err error
	if len(loginMsg.PostParams) != 0 {
		req, err = genPostParamReq(&loginMsg, &loginMsg.PostParams)
	} else {
		req, err = genGetParamReq(&loginMsg, &loginMsg.GetParams)
	}
	if logger.ErrHandle(err) {
		return nil
	}

	_, err = client.Do(req)
	if logger.ErrHandle(err) {
		return nil
	}
	return client4login.Jar
}

func createGetReq(url string, ref string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Himawari")
	if ref != "" {
		req.Header.Set("Referer", ref)
	}
	return
}

func createPostReq(url string, ref string, p url.Values) (req *http.Request, err error) {
	req, err = http.NewRequest("POST", url, strings.NewReader(p.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", ref)
	return
}

func genGetParamReq(j *entity.JsonMessage, gp *url.Values) (req *http.Request, err error) {
	req, err = createGetReq(j.URL, j.Referer)
	if logger.ErrHandle(err) {
		return
	}
	req.URL.RawQuery = gp.Encode()
	return
}

func genPostParamReq(j *entity.JsonMessage, pp *url.Values) (req *http.Request, err error) {
	req, err = createPostReq(j.URL, j.Referer, *pp)
	if logger.ErrHandle(err) {
		return
	}
	req.URL.RawQuery = j.GetParams.Encode()
	return
}
