package scanner

import (
	"Himawari/models/entity"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

//リダイレクト発生時、第３引数が元のリクエスト
func timeBasedAttack(s SendStruct, req []*http.Request) {

	//len(req)-1はリダイレクトがあったら元のほう

	var reqd []byte
	if len(req) == 1 {
		reqd, _ = httputil.DumpRequestOut(req[0], true)
		s.originalReq = reqd
	} else {
		reqd = s.originalReq
	}

	start := time.Now()
	resp, _ := client.Do(req[len(req)-1])
	end := time.Now()

	if compareAccessTime(s.jsonMessage.Time, (end.Sub(start)).Seconds(), s.kind) {

		respd, _ := httputil.DumpResponse(resp, true)

		newIssue := entity.Issue{
			URL:       s.jsonMessage.URL,
			Parameter: s.parameter,
			Kind:      s.kind,
			Getparam:  req[0].URL.Query(),
			Postparam: req[0].PostForm,
			Request:   string(reqd),
			Response:  string(respd),
		}
		*s.eachVulnIssue = append(*s.eachVulnIssue, newIssue)
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
		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq, _ = http.NewRequest(req[len(req)-1].Method, redirect.String(), strings.NewReader(req[len(req)-1].PostForm.Encode()))
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq, _ = http.NewRequest("GET", redirect.String(), nil)
			}
		} else {
			entity.Item.AppendItem(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		//s要検討(リダイレクト先のtimeと比較するのは難しい)
		timeBasedAttack(s, req)
	}
}
