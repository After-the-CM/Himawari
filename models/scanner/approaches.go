package scanner

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"Himawari/models/entity"
)

//リダイレクト発生時、第３引数が元のリクエスト
func timeBasedAttack(d determinant, req []*http.Request) {
	if len(req) == 1 {
		d.originalReq, _ = httputil.DumpRequestOut(req[0], true)
	}

	start := time.Now()
	resp, _ := client.Do(req[len(req)-1])
	end := time.Now()

	if compareAccessTime(d.jsonMessage.Time, (end.Sub(start)).Seconds(), d.kind) {

		dumpedResp, _ := httputil.DumpResponse(resp, true)

		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Parameter: d.parameter,
			Kind:      d.kind,
			Getparam:  req[0].URL.Query(),
			Postparam: req[0].PostForm,
			Request:   string(d.originalReq),
			Response:  string(dumpedResp),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
		entity.WholeIssue = append(entity.WholeIssue, newIssue)
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	//リダイレクト発生時
	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
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
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		//d要検討(リダイレクト先のtimeと比較するのは難しい)
		timeBasedAttack(d, req)
	}
}
