package scanner

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"reflect"
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
	candidate     *[]entity.JsonMessage
	randmark      string
}

const (
	settingTime   = 3
	tolerance     = 0.5
	osci          = "OS_Command_Injection"
	dirTraversal  = "Directory_Traversal"
	timeBasedSQLi = "Time_based_SQL_Injection"
	errBasedSQLi  = "Error_Based_SQL_Injection"
	reflectedXSS  = "Reflected_XSS"
	storedXSS     = "Stored_XSS"
	openRedirect  = "Open_Redirect"
	dirListing    = "Directory_Listing"
	httpHeaderi   = "HTTP_Header_Injection"
	csrf          = "Cross_Site_Request_Forgery"
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

var genRandmark = initRandmark(0)

//sleep時間は3秒で実行。誤差を考えるなら2.5秒くらい？

func compareAccessTime(originalTime float64, respTime float64, kind string) bool {
	if (respTime - originalTime) >= (settingTime - tolerance) {
		fmt.Println(kind)
		return true
	}
	return false
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
		if d.kind == openRedirect {
			d.setKeyValues(k, (payload), false, "GET")
		} else {
			d.setKeyValues(k, (v[0] + payload), false, "GET")
		}
	}

	//paramにpayload=1を追加する
	d.setKeyValues("Added by Himawari", payload, true, "POST")

	for k, v := range d.jsonMessage.PostParams {
		if d.kind == openRedirect {
			d.setKeyValues(k, (payload), false, "POST")
		} else {
			d.setKeyValues(k, (v[0] + payload), false, "POST")
		}
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

			req, err := genGetParamReq(d.jsonMessage, tmpUrlValues)
			if logger.ErrHandle(err) {
				return
			}

			d.approach(d, []*http.Request{req})

		case "POST":
			tmpUrlValues := copyUrlValues(&d.jsonMessage.PostParams)
			if addparam {
				tmpUrlValues.Add(payload, "1")
			} else {
				tmpUrlValues.Del(key)
				tmpUrlValues.Set(key, payload)
			}

			req, err := genPostParamReq(d.jsonMessage, tmpUrlValues)
			if logger.ErrHandle(err) {
				return
			}

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
		getPtReq, err := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		if logger.ErrHandle(err) {
			return
		}
		getPtReq.URL.Path = getPtReq.URL.Path + payload
		getPtReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{getPtReq})
	}

}

func (d determinant) setGetHeader(payload string) {
	//Header User-Agent
	d.parameter = "User-Agent"
	if !d.isAlreadyDetected() {
		getUAReq, err := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		if logger.ErrHandle(err) {
			return
		}
		getUAReq.Header.Set("User-Agent", getUAReq.UserAgent()+payload)
		getUAReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{getUAReq})
	}
	//Header Referer
	d.parameter = "Referer"
	if !d.isAlreadyDetected() {
		getRfReq, err := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		if logger.ErrHandle(err) {
			return
		}
		getRfReq.Header.Set("Referer", getRfReq.Referer()+payload)
		getRfReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{getRfReq})
	}

}

func (d determinant) setPostHeader(payload string) {
	//Header User-Agent
	d.parameter = "User-Agent"
	if !d.isAlreadyDetected() {
		postUAReq, err := createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
		if logger.ErrHandle(err) {
			return
		}
		postUAReq.PostForm = d.jsonMessage.PostParams
		postUAReq.Header.Set("User-Agent", postUAReq.UserAgent()+payload)
		postUAReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{postUAReq})
	}

	//Header Referer
	d.parameter = "Referer"
	if !d.isAlreadyDetected() {
		postRfReq, err := createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
		if logger.ErrHandle(err) {
			return
		}
		postRfReq.PostForm = d.jsonMessage.PostParams
		postRfReq.Header.Set("Referer", postRfReq.Referer()+payload)
		postRfReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{postRfReq})
	}
}

//fileにストリーム開く用
func readfile(fn string) *os.File {
	file, err := os.Open(fn)
	if logger.ErrHandle(err) {
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

func initRandmark(n int) func() string {
	cnt := n
	return func() string {
		cnt++
		return "Himawari" + fmt.Sprintf("%05d", cnt)
	}
}

func (d *determinant) gatherCandidates(j *entity.JsonNode) {
	//for _, v := range j.Messages {
	for i := 0; len(j.Messages) > i; i++ {

		d.randmark = genRandmark()
		d.setGetParam(d.randmark)
		d.randmark = genRandmark()
		d.setPostParam(d.randmark)

		//if fullscan{}
		/*
			if len(v.PostParams) != 0 {
				d.randmark = genRandmark()
				d.setPostUA(d.randmark)
				d.randmark = genRandmark()
				d.setPostRef(d.randmark)
			} else {
				d.randmark = genRandmark()
				d.setGetUA(d.randmark)
				d.randmark = genRandmark()
				d.setGetRef(d.randmark)
			}
		*/
	}

	for _, v := range j.Children {
		d.gatherCandidates(&v)
	}
}

// candidateの収集を行う
func (d *determinant) patrol(j entity.JsonNode, randmark string) {
	for _, v := range j.Messages {
		var req *http.Request
		var err error
		if len(v.PostParams) != 0 {
			req, err = genPostParamReq(&v, &v.PostParams)
		} else {
			req, err = genGetParamReq(&v, &v.GetParams)
		}

		if logger.ErrHandle(err) {
			return
		}

		resp, err := client.Do(req)
		if logger.ErrHandle(err) {
			return
		}
		body, err := io.ReadAll(resp.Body)
		if logger.ErrHandle(err) {
			return
		}
		targetResp := string(body)
		resp.Body.Close()

		if strings.Contains(targetResp, randmark) {
			if !isExist(d.candidate, v) {
				*d.candidate = append(*d.candidate, v)
			}
		}

		//通常のredirectならcrawl時に発見できているはず
	}
	for _, v := range j.Children {
		d.patrol(v, randmark)
	}
}

func (d determinant) setGetParam(payload string) {
	//paramにpayload=1を追加する
	//nameがない場合に追加するもの。nameの値を要件等
	d.setKeyValues("Added by Himawari", payload, true, "GET")

	for k, v := range d.jsonMessage.GetParams {
		d.setKeyValues(k, (v[0] + payload), false, "GET")
	}
}

func (d determinant) setPostParam(payload string) {
	//paramにpayload=1を追加する
	d.setKeyValues("Added by Himawari", payload, true, "POST")

	for k, v := range d.jsonMessage.PostParams {
		d.setKeyValues(k, (v[0] + payload), false, "POST")
	}
}

func (d determinant) setGetUA(payload string) {
	//Header User-Agent
	d.parameter = "User-Agent"
	if !d.isAlreadyDetected() {
		getUAReq, err := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		if logger.ErrHandle(err) {
			return
		}
		getUAReq.Header.Set("User-Agent", getUAReq.UserAgent()+payload)
		getUAReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{getUAReq})
	}
}

func (d determinant) setGetRef(payload string) {
	//Header Referer
	d.parameter = "Referer"
	if !d.isAlreadyDetected() {
		getRfReq, err := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
		if logger.ErrHandle(err) {
			return
		}
		getRfReq.Header.Set("Referer", getRfReq.Referer()+payload)
		getRfReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{getRfReq})
	}
}

func (d determinant) setPostUA(payload string) {
	//Header User-Agent
	d.parameter = "User-Agent"
	if !d.isAlreadyDetected() {
		postUAReq, err := createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
		if logger.ErrHandle(err) {
			return
		}
		postUAReq.PostForm = d.jsonMessage.PostParams
		postUAReq.Header.Set("User-Agent", postUAReq.UserAgent()+payload)
		postUAReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{postUAReq})
	}
}

func (d determinant) setPostRef(payload string) {
	//Header Referer
	d.parameter = "Referer"
	if !d.isAlreadyDetected() {
		var postRfReq *http.Request
		var err error
		if d.kind == csrf {
			postRfReq, err = createPostReq(d.jsonMessage.URL, "", d.jsonMessage.PostParams)
		} else {
			postRfReq, err = createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
		}

		if logger.ErrHandle(err) {
			return
		}

		postRfReq.PostForm = d.jsonMessage.PostParams
		postRfReq.Header.Set("Referer", postRfReq.Referer()+payload)
		postRfReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{postRfReq})
	}
}

func isExist(candidates *[]entity.JsonMessage, v entity.JsonMessage) bool {
	for _, candidate := range *candidates {
		if reflect.DeepEqual(candidate, v) {
			return true
		}
	}
	return false
}