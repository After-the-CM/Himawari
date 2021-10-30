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

	"github.com/PuerkitoBio/goquery"
)

//リダイレクト発生時req[0]がオリジナルのリクエスト
func timeBasedAttack(d determinant, req []*http.Request) {
	if len(req) == 1 {
		d.originalReq, _ = httputil.DumpRequestOut(req[0], true)
	}

	start := time.Now()
	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
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
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
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
	if len(req) == 1 {
		d.originalReq, _ = httputil.DumpRequestOut(req[0], true)
	}

	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var messages []string
	m := readfile("models/scanner/messages/" + d.kind + ".txt")
	msg := bufio.NewScanner(m)
	for msg.Scan() {
		messages = append(messages, msg.Text())
	}
	dumpedResp, _ := httputil.DumpResponse(resp, true)

	body, _ := io.ReadAll(resp.Body)
	targetResp := string(body)

	var u string
	if d.kind == DirList {
		u = req[0].URL.String()
	} else {
		u = d.jsonMessage.URL
	}

	for _, msg := range messages {
		if strings.Contains(targetResp, msg) {
			fmt.Println(d.kind)
			newIssue := entity.Issue{
				URL:       u,
				Parameter: d.parameter,
				Kind:      d.kind,
				Getparam:  req[0].URL.Query(),
				Postparam: req[0].PostForm,
				Request:   string(d.originalReq),
				Response:  string(dumpedResp),
			}
			*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
			entity.WholeIssue = append(entity.WholeIssue, newIssue)
			break
		}
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, _ := url.Parse(location)
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
		} else {
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		stringMatching(d, req)
	}
}

func detectReflectedXSS(d determinant, req []*http.Request) {
	if len(req) == 1 {
		d.originalReq, _ = httputil.DumpRequestOut(req[0], true)
	}

	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	dumpedResp, _ := httputil.DumpResponse(resp, true)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		injectedPayload := s.Text()
		if strings.Contains(injectedPayload, "alert(1)") {
			fmt.Println(d.kind)
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
			return false
		}
		return true
	})

	io.ReadAll(resp.Body)
	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, _ := url.Parse(location)
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
		} else {
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		detectReflectedXSS(d, req)
	}
}

func detectStoredXSS(d determinant, req []*http.Request) {
	if len(req) == 1 {
		d.originalReq, _ = httputil.DumpRequestOut(req[0], true)
	}

	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var dumpedResp []byte
	b := false
	for _, v := range *d.candidate {
		var inspectReq *http.Request
		if len(v.PostParams) != 0 {
			inspectReq = genPostParamReq(&v, &v.PostParams)
		} else {
			inspectReq = genGetParamReq(&v, &v.GetParams)
		}

		inspectResp, err := client.Do(inspectReq)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		dumpedResp, _ = httputil.DumpResponse(inspectResp, true)

		doc, err := goquery.NewDocumentFromReader(inspectResp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			injectedPayload := s.Text()
			if strings.Contains(injectedPayload, "alert(\""+d.randmark+"\")") {
				fmt.Println(d.kind)
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
				b = true
				return false
			}
			return true
		})

		io.ReadAll(inspectResp.Body)
		inspectResp.Body.Close()

		if b {
			return
		}
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()
}

func searchRandmark(d determinant, req []*http.Request) {
	resp, err := client.Do(req[len(req)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	targetResp := string(body)
	io.ReadAll(resp.Body)
	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, _ := url.Parse(location)
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
		} else {
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		searchRandmark(d, req)
	}

	if strings.Contains(targetResp, d.randmark) {
		// reflect
		return
	} else {
		// stored
		d.patrol(entity.JsonNodes, d.randmark)
	}
}
