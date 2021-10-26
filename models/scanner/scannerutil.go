package scanner

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

type determinant struct {
	jsonMessage   *entity.JsonMessage
	parameter     string
	kind          string
	originalReq   []byte
	approach      func(d determinant, req []*http.Request)
	eachVulnIssue *[]entity.Issue
}

const (
	PayloadTime = 3
	tolerance   = 0.5
	OSCI        = "OS_Command_Injection"
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

//sleep時間は3秒で実行。誤差を考えるなら2.5秒くらい？

func compareAccessTime(originalTime float64, respTime float64, kind string) bool {
	if (respTime - originalTime) >= (PayloadTime - tolerance) {
		fmt.Fprintln(os.Stderr, kind)
		return true
	}
	return false
}

func createGetReq(url string, ref string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", ref)
	return req

}

func createPostReq(url string, ref string, p url.Values) *http.Request {
	req, _ := http.NewRequest("POST", url, strings.NewReader(p.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", ref)
	return req

}

//jsonMessageのissueに同じパラメーターで、同じ種類の脆弱性があるか確認する
func (d determinant) isAlreadyDetected() bool {
	//ワンちゃん一番最後だけでええんちゃう？
	for _, v := range *d.eachVulnIssue {
		if v.Parameter == d.parameter && v.Kind == d.kind && v.URL == d.jsonMessage.URL {
			return true
		}
	}
	return false
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

func genGetHeaderReq(req *http.Request, gp *url.Values) *http.Request {
	req.URL.RawQuery = gp.Encode()
	return req
}

func genPostHeaderReq(req *http.Request, gp *url.Values) *http.Request {
	req.URL.RawQuery = gp.Encode()
	return req
}

//jsonMessageのparamをコピー
func copyUrlValues(u *url.Values) *url.Values {
	tmp := url.Values{}

	for k, v := range *u {
		tmp[k] = v
	}
	return &tmp
}

func (d determinant) setParam(payload string) {
	//paramにpayload=1を追加する
	//nameがない場合に追加するもの。nameの値を要件等
	d.setKeyValues("Added by Himawari", payload, true, "GET")

	for k, v := range d.jsonMessage.GetParams {
		d.setKeyValues(k, (v[0] + payload), false, "GET")
	}

	//paramにpayload=1を追加する
	d.setKeyValues("Added by Himawari", payload, true, "POST")

	for k, v := range d.jsonMessage.PostParams {
		d.setKeyValues(k, (v[0] + payload), false, "POST")
	}
}

func (d determinant) setKeyValues(key string, payload string, addparam bool, method string) {
	d.parameter = key

	if !d.isAlreadyDetected() {
		switch method {
		case "GET":
			tmpUrlValues := copyUrlValues(&d.jsonMessage.GetParams)
			if addparam {
				tmpUrlValues.Add(payload, "1")
			} else {
				tmpUrlValues.Del(key)
				tmpUrlValues.Set(key, payload)
			}

			req := genGetParamReq(d.jsonMessage, tmpUrlValues)
			d.approach(d, []*http.Request{req})

		case "POST":
			tmpUrlValues := copyUrlValues(&d.jsonMessage.PostParams)
			if addparam {
				tmpUrlValues.Add(payload, "1")
			} else {
				tmpUrlValues.Del(key)
				tmpUrlValues.Set(key, payload)
			}

			req := genPostParamReq(d.jsonMessage, tmpUrlValues)
			req.PostForm = *tmpUrlValues
			d.approach(d, []*http.Request{req})
		default:
			fmt.Fprintf(os.Stderr, "No support method\n")
		}
	}
}

func (d determinant) setHeaderDocumentRoot(payload string) {
	d.parameter = "Path"
	if !d.isAlreadyDetected() {
		getPtReq := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		getPtReq.URL.Path = getPtReq.URL.Path + payload

		req := genGetHeaderReq(getPtReq, &d.jsonMessage.GetParams)
		d.approach(d, []*http.Request{req})
	}

}

func (d determinant) setGetHeader(payload string) {
	//Header User-Agent
	d.parameter = "User-Agent"
	if !d.isAlreadyDetected() {
		getUAReq := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		getUAReq.Header.Set("User-Agent", getUAReq.UserAgent()+payload)

		req := genGetHeaderReq(getUAReq, &d.jsonMessage.GetParams)
		d.approach(d, []*http.Request{req})
	}
	//Header Referer
	d.parameter = "Referer"
	if !d.isAlreadyDetected() {
		getRfReq := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		getRfReq.Header.Set("Referer", getRfReq.Referer()+payload)

		req := genGetHeaderReq(getRfReq, &d.jsonMessage.GetParams)
		d.approach(d, []*http.Request{req})
	}

}

func (d determinant) setPostHeader(payload string) {
	//Header User-Agent
	d.parameter = "User-Agent"
	if !d.isAlreadyDetected() {
		postUAReq := createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
		postUAReq.PostForm = d.jsonMessage.PostParams
		postUAReq.Header.Set("User-Agent", postUAReq.UserAgent()+payload)

		req := genPostHeaderReq(postUAReq, &d.jsonMessage.GetParams)
		d.approach(d, []*http.Request{req})
	}

	//Header Referer
	d.parameter = "Referer"
	if !d.isAlreadyDetected() {
		postRfReq := createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
		postRfReq.PostForm = d.jsonMessage.PostParams
		postRfReq.Header.Set("Referer", postRfReq.Referer()+payload)

		req := genPostHeaderReq(postRfReq, &d.jsonMessage.GetParams)
		d.approach(d, []*http.Request{req})
	}
}

//fileにストリーム開く用
func readfile(fn string) *os.File {
	file, err := os.Open(fn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return file
}

func retrieveJsonMessage(j *entity.JsonNode) *entity.JsonMessage {
	if len(j.Messages) != 0 {
		return &j.Messages[0]
	}
	for _, v := range j.Children {
		if len(v.Messages) != 0 {
			return retrieveJsonMessage(&v)
		}
	}
	return nil
}
