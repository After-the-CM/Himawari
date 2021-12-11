package scanner

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

type determinant struct {
	jsonMessage   *entity.JsonMessage
	parameter     string
	payload       string
	kind          string
	originalReq   []byte
	approach      func(d determinant, req []*http.Request)
	eachVulnIssue *[]entity.Issue
	landmark      string
	cookie        entity.JsonCookie
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

var QuickScan bool = false

var genLandmark = initLandmark(0)

func SetGenLandmark(n int) {
	genLandmark = initLandmark(n)
}

//sleep時間は3秒で実行。誤差を考えるなら2.5秒くらい？
func compareAccessTime(originalTime float64, respTime float64) bool {
	return (respTime - originalTime) >= (settingTime - tolerance)
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

func (d determinant) prepareLandmark(payload string) {
	d.landmark = genLandmark()
	d.setKeyValues("Added by Himawari", strings.Replace(payload, "[landmark]", d.landmark, 1), true, "GET")

	for key, value := range d.jsonMessage.GetParams {
		d.landmark = genLandmark()
		d.setKeyValues(key, (value[0] + strings.Replace(payload, "[landmark]", d.landmark, 1)), false, "GET")
	}

	d.landmark = genLandmark()
	d.setKeyValues("Added by Himawari", strings.Replace(payload, "[landmark]", d.landmark, 1), true, "POST")

	for key, value := range d.jsonMessage.PostParams {
		d.landmark = genLandmark()
		d.setKeyValues(key, (value[0] + strings.Replace(payload, "[landmark]", d.landmark, 1)), false, "POST")
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

			d.payload = payload
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
			d.payload = payload
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
		d.payload = getPtReq.URL.Path + payload
		getPtReq.URL.Path = getPtReq.URL.Path + payload
		getPtReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{getPtReq})
	}

}

func (d determinant) setCookie(cookie entity.JsonCookie, payload string) {
	d.parameter = cookie.Name
	payload = strings.Replace(payload, " ", "%20", -1)
	d.payload = cookie.Value + payload

	if !d.isAlreadyDetected() {
		d.cookie = entity.JsonCookie{
			Path:  cookie.Path,
			Name:  cookie.Name,
			Value: cookie.Value + payload,
		}
		if len(d.jsonMessage.PostParams) == 0 {
			req, err := createGetReq(d.jsonMessage.URL, d.jsonMessage.Referer)
			if logger.ErrHandle(err) {
				return
			}
			req.URL.RawQuery = d.jsonMessage.GetParams.Encode()

			req.AddCookie(&http.Cookie{
				Path:  d.cookie.Path,
				Name:  d.cookie.Name,
				Value: d.cookie.Value,
			})
			d.approach(d, []*http.Request{req})
		} else {
			req, err := createPostReq(d.jsonMessage.URL, d.jsonMessage.Referer, d.jsonMessage.PostParams)
			if logger.ErrHandle(err) {
				return
			}
			req.PostForm = d.jsonMessage.PostParams
			req.URL.RawQuery = d.jsonMessage.GetParams.Encode()

			req.AddCookie(&http.Cookie{
				Path:  d.cookie.Path,
				Name:  d.cookie.Name,
				Value: d.cookie.Value,
			})

			d.approach(d, []*http.Request{req})
		}
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
		d.payload = getUAReq.UserAgent() + payload
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
		d.payload = getRfReq.Referer() + payload
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
		d.payload = postUAReq.UserAgent() + payload
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
		d.payload = postRfReq.Referer() + payload
		postRfReq.Header.Set("Referer", postRfReq.Referer()+payload)
		postRfReq.URL.RawQuery = d.jsonMessage.GetParams.Encode()

		d.approach(d, []*http.Request{postRfReq})
	}
}

var exe, _ = os.Executable()
var dir = filepath.Dir(exe)

//fileにストリーム開く用
func readfile(fn string) *os.File {
	file, err := os.Open(dir + "/" + fn)
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

func initLandmark(n int) func() string {
	cnt := n
	return func() string {
		cnt++
		return fmt.Sprintf("65535%05d", cnt)
	}
}

func (d *determinant) gatherCandidates() {
	d.prepareLandmark("[landmark]")

	if len(d.jsonMessage.PostParams) != 0 {
		d.landmark = genLandmark()
		d.setPostUA(d.landmark)
		d.landmark = genLandmark()
		d.setPostRef(d.landmark)
	} else {
		d.landmark = genLandmark()
		d.setGetUA(d.landmark)
		d.landmark = genLandmark()
		d.setGetRef(d.landmark)
	}
}

// candidateの収集を行う
func (d *determinant) patrol(j entity.JsonNode, landmark string) {
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

		time.Sleep(entity.RequestDelay)

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

		if strings.Contains(targetResp, landmark) {
			if !isExist(&d.jsonMessage.Candidate, v) {
				d.jsonMessage.Candidate = append(d.jsonMessage.Candidate, v)
			}
		}

		//通常のredirectならcrawl時に発見できているはず
	}
	for _, v := range j.Children {
		d.patrol(v, landmark)
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
		d.payload = getUAReq.UserAgent() + payload
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
		d.payload = getRfReq.Referer() + payload
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
		d.payload = postUAReq.UserAgent() + payload
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
		d.payload = postRfReq.Referer() + payload
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

func (d determinant) extractCookie(cookies []*http.Cookie) []*http.Cookie {
	cookieExtract := make([]*http.Cookie, 0)
	for _, cookie := range cookies {
		if cookie.Name != d.cookie.Name {
			cookieExtract = append(cookieExtract, cookie)
		}
	}
	return cookieExtract
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

	time.Sleep(entity.RequestDelay)

	_, err = client.Do(req)
	if logger.ErrHandle(err) {
		return nil
	}
	return client4login.Jar
}

func CalcApproximateTime(j *entity.JsonNode) {
	accessTime := retrieveAccessTime(j)
	fmt.Println("アクセス時間:", accessTime)
	if entity.RequestDelay.Microseconds() != 0 {
		fmt.Println("遅延", entity.RequestDelay)
	}
	accessTime += float64(entity.RequestDelay.Milliseconds())

	msgNum := countMsg(j)
	paramNum := countParam(j)
	cookieNum := countCookie(j)
	accessNum := ((msgNum * 5) + paramNum + cookieNum) * 315
	fmt.Println("メッセージ数:", msgNum)
	fmt.Println("パラメータ数:", paramNum)
	fmt.Println("Cookie数:", cookieNum)
	fmt.Println("\tアクセス数:", accessNum)
	if !QuickScan {
		fmt.Println("full scan。アクセス数が増加します。")
		accessNum += ((msgNum * 2) + paramNum + cookieNum) * msgNum
		fmt.Println("\tアクセス数:", accessNum)
	}

	accessTime += float64(accessNum) * 0.0001 // 処理時間としてマージン

	if loginMsg.URL != "" {
		fmt.Println("login有。アクセス数が増加します。")
		accessNum *= 2
		fmt.Println("\tアクセス数:", accessNum)
	}

	fmt.Println("予想診断時間は", int(accessTime*float64(accessNum))/60000, "分です。")
	fmt.Println("※診断時間はあくまで目安です。診断対象の挙動によって増減します。")
}

func countMsg(j *entity.JsonNode) (msgNum int) {
	msgNum = len(j.Messages)
	for i := range j.Children {
		msgNum += countMsg(&j.Children[i])
	}
	return msgNum
}

func countCookie(j *entity.JsonNode) (cookieNum int) {
	cookieNum += len(j.Cookies)
	for i := range j.Children {
		cookieNum += countMsg(&j.Children[i])
	}
	return cookieNum
}

func countParam(j *entity.JsonNode) (paramNum int) {
	for i := range j.Messages {
		paramNum += len(j.Messages[i].GetParams)
		paramNum += len(j.Messages[i].PostParams)
	}
	for i := range j.Children {
		paramNum += countMsg(&j.Children[i])
	}
	return paramNum
}

func retrieveAccessTime(j *entity.JsonNode) float64 {
	var sum, cnt, time float64
	for i := range j.Messages {
		sum += j.Messages[i].Time
		cnt++
	}
	time = sum / cnt

	for i := range j.Children {
		time = retrieveAccessTime(&j.Children[i])
	}
	return time
}
