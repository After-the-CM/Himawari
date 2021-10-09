package scanner

import (
	"Himawari/models/entity"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type SendStruct struct {
	jsonMessage *entity.JsonMessage
	//req         *http.Request
	parameter     string
	kind          string
	sendMethod    func(s SendStruct, req []*http.Request)
	eachVulnIssue *[]entity.Issue
}

const (
	PayloadTime = 3
	MarginTime  = 0.5
	OSCI        = "OS Command Injection"
)

//sleep時間は3秒で実行。誤差を考えるなら2.5秒くらい？

func compareAccessTime(origin float64, resp float64) bool {

	if (resp-origin) >= (PayloadTime-MarginTime) && (PayloadTime+MarginTime) >= (resp-origin) {
		fmt.Fprintln(os.Stderr, "os command injection!!")
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
func (s SendStruct) isDetectedIssue() bool {

	//ワンちゃん一番最後だけでええんちゃう？
	//	if j.issue != nil {
	for _, v := range *s.eachVulnIssue {
		if v.Parameter == s.parameter && v.Kind == s.kind && v.URL == s.jsonMessage.URL {
			return true
		}

	}
	/*
		} else {
			j.issue = &[]Foundpage{}
		}
	*/
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

func extractPostValues(req *http.Request) url.Values {
	p := req.ParseForm()
	fmt.Fprintln(io.Discard, p)
	return req.PostForm
}

//リダイレクト発生時、第３引数が元のリクエスト
func timeBasedAttack(s SendStruct, req []*http.Request) {

	//client := new(http.Client)

	//client.Do(req)をする前に実行しないとリクエスト内容が消えてしまう。
	//len(req)-1はリダイレクトがあったら元のほう
	reqd, _ := httputil.DumpRequestOut(req[0], true)

	start := time.Now()
	resp, _ := client.Do(req[0])
	end := time.Now()

	if compareAccessTime(s.jsonMessage.Time, (end.Sub(start)).Seconds()) {

		/*
			if req.Body != nil {
				req.Body, _ = req.GetBody()
			}
		*/

		//	fmt.Println(string(req))

		respd, _ := httputil.DumpResponse(resp, true)

		newIssue := entity.Issue{
			URL: s.jsonMessage.URL,
			//URL:       req.URL.String(),
			Parameter: s.parameter,
			Kind:      s.kind,
			Getparam:  req[0].URL.Query(),
			//Postparam: extractPostValues(req),
			Postparam: req[0].PostForm,
			Request:   string(reqd),
			Response:  string(respd),
		}
		//osci内でのグローバル変数の
		//nodeIssue = append(nodeIssue, s)
		*s.eachVulnIssue = append(*s.eachVulnIssue, newIssue)
		//サイト全体のIssues
		//Issues = append(Issues, s)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Fprintln(io.Discard, string(body))
	resp.Body.Close()

	//リダイレクト発生時
	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		//307想定、動くなら
		l, _ := url.Parse(location)
		redirect := req[len(req)-1].URL.ResolveReference(l)
		if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
			redirectReq, _ = http.NewRequest(req[len(req)-1].Method, redirect.String(), strings.NewReader(req[len(req)-1].PostForm.Encode()))
			redirectReq.PostForm = req[len(req)-1].PostForm
		} else {
			redirectReq, _ = http.NewRequest("GET", redirect.String(), nil)
		}
		req = append(req, redirectReq)
		//s要検討(リダイレクト先のtimeと比較するのは難しい)
		timeBasedAttack(s, req)
	}
}

//jsonMessageのparamをコピー
func copyUrlValues(u *url.Values) *url.Values {
	tmp := url.Values{}

	for k, v := range *u {
		tmp[k] = v
	}

	return &tmp
}

func (s SendStruct) setKeyValues(key string, payload string, addparam bool, method string) {
	s.parameter = key

	if !s.isDetectedIssue() {
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
			s.sendMethod(s, []*http.Request{req})

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
			s.sendMethod(s, []*http.Request{req})
		default:
			fmt.Fprintf(os.Stderr, "No support method\n")

		}
	}
}

func (s SendStruct) setParam(payload string) {
	//paramにpayload=1を追加する
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

func (s SendStruct) setHeaderDocumentRoot(payload string) {

	s.parameter = "Path"
	if !s.isDetectedIssue() {
		getPtReq := createGetReq(s.jsonMessage)
		getPtReq.URL.Path = getPtReq.URL.Path + payload

		req := genGetHeaderReq(getPtReq, "Path", &s.jsonMessage.GetParams)
		s.sendMethod(s, []*http.Request{req})
	}

}

func (s SendStruct) setGetHeader(payload string) {

	//Header User-Agent
	s.parameter = "User-Agent"
	if !s.isDetectedIssue() {
		getUAReq := createGetReq(s.jsonMessage)
		getUAReq.Header.Set("User-Agent", getUAReq.UserAgent()+payload)

		req := genGetHeaderReq(getUAReq, "User-Agent", &s.jsonMessage.GetParams)
		s.sendMethod(s, []*http.Request{req})
	}
	//Header Referer
	s.parameter = "Referer"
	if !s.isDetectedIssue() {
		getRfReq := createGetReq(s.jsonMessage)
		getRfReq.Header.Set("Referer", getRfReq.Referer()+payload)

		req := genGetHeaderReq(getRfReq, "Referer", &s.jsonMessage.GetParams)
		s.sendMethod(s, []*http.Request{req})
	}

	//改行の文字コードも追加できる
	//xxxx.Header.Add("test", "%0A")

	//Header Method
	/*
		req := j.createreq()
		req.Method = "GET%3Bls"
		j.getRequestOfHeader(req,"Header Method")
	*/

}

func (s SendStruct) setPostHeader(payload string) {
	//Header User-Agent
	s.parameter = "User-Agent"
	if !s.isDetectedIssue() {
		postUAReq := createPostReq(s.jsonMessage, s.jsonMessage.PostParams)
		postUAReq.PostForm = s.jsonMessage.PostParams
		postUAReq.Header.Set("User-Agent", postUAReq.UserAgent()+payload)

		req := genPostHeaderReq(postUAReq, s.parameter, &s.jsonMessage.GetParams)
		s.sendMethod(s, []*http.Request{req})
	}

	//Header Referer
	s.parameter = "Referer"
	if !s.isDetectedIssue() {
		postRfReq := createPostReq(s.jsonMessage, s.jsonMessage.PostParams)
		postRfReq.PostForm = s.jsonMessage.PostParams
		postRfReq.Header.Set("Referer", postRfReq.Referer()+payload)

		req := genPostHeaderReq(postRfReq, s.parameter, &s.jsonMessage.GetParams)
		s.sendMethod(s, []*http.Request{req})
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
