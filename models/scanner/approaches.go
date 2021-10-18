package scanner

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"Himawari/models/entity"
)

//リダイレクト発生時、第３引数が元のリクエスト
func timeBasedAttack(d determinant, req []*http.Request) {
	var reqd []byte
	if len(req) == 1 {
		reqd, _ = httputil.DumpRequestOut(req[0], true)
		d.originalReq = reqd
	} else {
		reqd = d.originalReq
	}

	start := time.Now()
	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	end := time.Now()

	if compareAccessTime(d.jsonMessage.Time, (end.Sub(start)).Seconds(), d.kind) {

		respd, _ := httputil.DumpResponse(resp, true)

		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Parameter: d.parameter,
			Kind:      d.kind,
			Getparam:  req[0].URL.Query(),
			Postparam: req[0].PostForm,
			Request:   string(reqd),
			Response:  string(respd),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprintln(io.Discard, string(body))
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

func stringMatching(d determinant, req []*http.Request) {
	var reqd []byte
	if len(req) == 1 {
		reqd, _ = httputil.DumpRequestOut(req[0], true)
		d.originalReq = reqd
	} else {
		reqd = d.originalReq
	}
	//reqd, _ := httputil.DumpRequestOut(req[0], true)

	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var messages []string
	//messages := make([]string, ..., ...)

	// s.kindはconstで定義されている。
	// error message等の比較したい特定の文字列が入っているファイルを指定する。
	m := readfile("models/scanner/messages/" + string(d.kind) + ".txt")
	msg := bufio.NewScanner(m)
	for msg.Scan() {
		messages = append(messages, msg.Text())
	}
	respd, _ := httputil.DumpResponse(resp, true)

	body, _ := io.ReadAll(resp.Body)
	var targetResp string
	if d.kind == "HTTPHeaderInjection" || d.kind == "OpenRedirect" {
		// http.response.Headerはmap[string][]で単にstringに変換できない。
		// とりあえずReferer
		targetResp = resp.Header.Get("Referer")
	} else {
		targetResp = string(body)
	}

	for _, msg := range messages {
		if strings.Contains(targetResp, msg) {
			fmt.Println(d.kind)
			//respd, _ := httputil.DumpResponse(resp, true)

			newIssue := entity.Issue{
				URL: d.jsonMessage.URL,
				//URL:       req.URL.String(),
				Parameter: d.parameter,
				Kind:      d.kind,
				Getparam:  req[0].URL.Query(),
				//Postparam: extractPostValues(req),
				Postparam: req[0].PostForm,
				Request:   string(reqd),
				Response:  string(respd),
			}
			*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
			fmt.Fprintln(io.Discard, string(body))
			resp.Body.Close()
			return
		}
	}

	//リダイレクト発生時
	//location := resp.Header.Get("Location")
	//fmt.Fprintln(io.Discard, string(body))
	//resp.Body.Close()
	/*
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
			stringMatching(d, req)
		}
	*/
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
		stringMatching(d, req)
	}
}
