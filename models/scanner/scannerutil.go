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

type SendStruct struct {
	jsonMessage *entity.JsonMessage
	//req         *http.Request
	parameter     string
	kind          string
	approach      func(s SendStruct, req []*http.Request)
	eachVulnIssue *[]entity.Issue
}

const (
	PayloadTime = 3
	tolerance   = 0.5
	OSCI        = "OS Command Injection"
)

var jar, _ = cookiejar.New(nil)
var client = &http.Client{
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

//sleep時間は3秒で実行。誤差を考えるなら2.5秒くらい？

func compareAccessTime(originalTime float64, respTime float64, kind string) bool {

	if (respTime-originalTime) >= (PayloadTime-tolerance) && (PayloadTime+tolerance) >= (respTime-originalTime) {
		fmt.Fprintln(os.Stderr, kind)
		return true
	}

	return false
}

func createGetReq(j *entity.JsonMessage) *http.Request {

	req, _ := http.NewRequest("GET", j.URL, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", j.Referer)
	return req

}

func createPostReq(j *entity.JsonMessage, p url.Values) *http.Request {
	req, _ := http.NewRequest("POST", j.URL, strings.NewReader(p.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Himawari")
	req.Header.Set("Referer", j.Referer)
	return req

}

//jsonMessageのissueに同じパラメーターで、同じ種類の脆弱性があるか確認する
func (s SendStruct) isAlreadyDetected() bool {

	//ワンちゃん一番最後だけでええんちゃう？
	for _, v := range *s.eachVulnIssue {
		if v.Parameter == s.parameter && v.Kind == s.kind && v.URL == s.jsonMessage.URL {
			return true
		}

	}

	return false
}

//func genGetParamReq(j *JsonMessage, param string, kind string, gp *url.Values, s SendStruct) *http.Request {

func genGetParamReq(j *entity.JsonMessage, gp *url.Values) *http.Request {
	req := createGetReq(j)
	req.URL.RawQuery = gp.Encode()
	return req
}

func genPostParamReq(j *entity.JsonMessage, pp *url.Values) *http.Request {
	req := createPostReq(j, *pp)
	req.URL.RawQuery = j.GetParams.Encode()
	return req
}

func genGetHeaderReq(req *http.Request, param string, gp *url.Values) *http.Request {
	req.URL.RawQuery = gp.Encode()
	return req
}

func genPostHeaderReq(req *http.Request, param string, gp *url.Values) *http.Request {
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

func (s SendStruct) setParam(payload string) {
	//paramにpayload=1を追加する
	//nameがない場合に追加するもの。nameの値を要件等
	s.setKeyValues("Added by Himawari", payload, true, "GET")

	for k, v := range s.jsonMessage.GetParams {
		s.setKeyValues(k, (v[0] + payload), false, "GET")
	}

	//paramにpayload=1を追加する
	s.setKeyValues("Added by Himawari", payload, true, "POST")

	for k, v := range s.jsonMessage.PostParams {
		s.setKeyValues(k, (v[0] + payload), false, "POST")
	}
}
func (s SendStruct) setKeyValues(key string, payload string, addparam bool, method string) {
	s.parameter = key

	if !s.isAlreadyDetected() {
		switch method {
		case "GET":
			tmpUrlValues := copyUrlValues(&s.jsonMessage.GetParams)

			if addparam {
				tmpUrlValues.Add(payload, "1")
				//place = key
			} else {
				tmpUrlValues.Del(key)
				tmpUrlValues.Set(key, payload)
			}

			req := genGetParamReq(s.jsonMessage, tmpUrlValues)
			s.approach(s, []*http.Request{req})

		case "POST":

			tmpUrlValues := copyUrlValues(&s.jsonMessage.PostParams)
			if addparam {
				tmpUrlValues.Add(payload, "1")
			} else {
				tmpUrlValues.Del(key)
				tmpUrlValues.Set(key, payload)
			}

			req := genPostParamReq(s.jsonMessage, tmpUrlValues)
			req.PostForm = *tmpUrlValues
			s.approach(s, []*http.Request{req})
		default:
			fmt.Fprintf(os.Stderr, "No support method\n")

		}
	}
}

func (s SendStruct) setHeaderDocumentRoot(payload string) {

	s.parameter = "Path"
	if !s.isAlreadyDetected() {
		getPtReq := createGetReq(s.jsonMessage)
		getPtReq.URL.Path = getPtReq.URL.Path + payload

		req := genGetHeaderReq(getPtReq, "Path", &s.jsonMessage.GetParams)
		s.approach(s, []*http.Request{req})
	}

}

func (s SendStruct) setGetHeader(payload string) {

	//Header User-Agent
	s.parameter = "User-Agent"
	if !s.isAlreadyDetected() {
		getUAReq := createGetReq(s.jsonMessage)
		getUAReq.Header.Set("User-Agent", getUAReq.UserAgent()+payload)

		req := genGetHeaderReq(getUAReq, "User-Agent", &s.jsonMessage.GetParams)
		s.approach(s, []*http.Request{req})
	}
	//Header Referer
	s.parameter = "Referer"
	if !s.isAlreadyDetected() {
		getRfReq := createGetReq(s.jsonMessage)
		getRfReq.Header.Set("Referer", getRfReq.Referer()+payload)

		req := genGetHeaderReq(getRfReq, "Referer", &s.jsonMessage.GetParams)
		s.approach(s, []*http.Request{req})
	}

}

func (s SendStruct) setPostHeader(payload string) {
	//Header User-Agent
	s.parameter = "User-Agent"
	if !s.isAlreadyDetected() {
		postUAReq := createPostReq(s.jsonMessage, s.jsonMessage.PostParams)
		postUAReq.PostForm = s.jsonMessage.PostParams
		postUAReq.Header.Set("User-Agent", postUAReq.UserAgent()+payload)

		req := genPostHeaderReq(postUAReq, "User-Agent", &s.jsonMessage.GetParams)
		s.approach(s, []*http.Request{req})
	}

	//Header Referer
	s.parameter = "Referer"
	if !s.isAlreadyDetected() {
		postRfReq := createPostReq(s.jsonMessage, s.jsonMessage.PostParams)
		postRfReq.PostForm = s.jsonMessage.PostParams
		postRfReq.Header.Set("Referer", postRfReq.Referer()+payload)

		req := genPostHeaderReq(postRfReq, "Referer", &s.jsonMessage.GetParams)
		s.approach(s, []*http.Request{req})
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
