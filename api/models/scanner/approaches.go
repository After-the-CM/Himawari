package scanner

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"Himawari/models/entity"
	"Himawari/models/logger"

	"github.com/PuerkitoBio/goquery"
)

//リダイレクト発生時req[0]がオリジナルのリクエスト
func timeBasedAttack(d determinant, req []*http.Request) {
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	start := time.Now()
	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}
	end := time.Now()

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	if compareAccessTime(d.jsonMessage.Time, (end.Sub(start)).Seconds(), d.kind) {
		dumpedResp, err := httputil.DumpResponse(resp, true)

		//string()は引数がnilの場合でもnilぽエラーが出ない
		logger.ErrHandle(err)
		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Kind:      d.kind,
			Parameter: d.parameter,
			Payload:   d.payload,
			Evidence:  "Response delay: " + fmt.Sprint(end.Sub(start)),
			Request:   string(d.originalReq),
			Response:  string(dumpedResp),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
		entity.WholeIssue = append(entity.WholeIssue, newIssue)
		entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	//リダイレクト発生時
	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request

		//locationのParseに失敗した場合リダイレクトできないのでreturn
		l, err := url.Parse(location)
		if logger.ErrHandle(err) {
			return
		}
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				var err error
				redirectReq, err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
				if logger.ErrHandle(err) {
					return
				}
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq,err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
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
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	var messages []string
	m := readfile("models/scanner/messages/" + d.kind + ".txt")
	msg := bufio.NewScanner(m)
	for msg.Scan() {
		messages = append(messages, msg.Text())
	}

	//ここでdumpを行わないとResponseBodyが取れない。
	dumpedResp, err := httputil.DumpResponse(resp, true)
	logger.ErrHandle(err)

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	body, err := io.ReadAll(resp.Body)
	if !logger.ErrHandle(err) {

		targetResp := string(body)

		var u string
		if d.kind == dirListing {
			u = req[0].URL.String()
		} else {
			u = d.jsonMessage.URL
		}

		for _, msg := range messages {
			if strings.Contains(targetResp, msg) {
				fmt.Println(d.kind)
				newIssue := entity.Issue{
					URL:       u,
					Kind:      d.kind,
					Parameter: d.parameter,
					Payload:   d.payload,
					Evidence:  "Text match: " + msg,
					Request:   string(d.originalReq),
					Response:  string(dumpedResp),
				}
				*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
				entity.WholeIssue = append(entity.WholeIssue, newIssue)
				entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
				break
			}
		}
	}

	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, err := url.Parse(location)
		if logger.ErrHandle(err) {
			return
		}
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				var err error
				redirectReq, err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
				if logger.ErrHandle(err) {
					return
				}
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq,err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
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
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	//ここでdumpを行わないとResponseBodyが取れない。
	dumpedResp, err := httputil.DumpResponse(resp, true)
	logger.ErrHandle(err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	logger.ErrHandle(err)

	var flg bool
	var evidence string
	doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		injectedPayload := s.Text()
		if strings.Contains(injectedPayload, "alert(\""+d.landmark+"\")") {
			evidence = "alert(\"" + d.landmark + "\")"
			flg = true
			return false
		}
		return true
	})

	if !flg {
		doc.Find("*").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			href, _ := s.Attr("href")
			if strings.HasPrefix(href, "javascript:alert(\""+d.landmark+"\")") {
				evidence = "javascript:alert(\"" + d.landmark + "\")"
				flg = true
				return false
			}
			src, _ := s.Attr("src")
			if strings.HasPrefix(src, "javascript:alert(\""+d.landmark+"\")") {
				evidence = "javascript:alert(\"" + d.landmark + "\")"
				flg = true
				return false
			}
			if src == "x" {
				onerror, _ := s.Attr("onerror")
				if strings.Contains(onerror, "alert(\""+d.landmark+"\")") {
					evidence = "alert(\"" + d.landmark + "\")"
					flg = true
					return false
				}
			}
			onmouseover, _ := s.Attr("onmouseover")
			if strings.Contains(onmouseover, "alert(\""+d.landmark+"\")") {
				evidence = "alert(\"" + d.landmark + "\")"
				flg = true
				return false
			}
			return true
		})
	}

	if flg {
		fmt.Println(d.kind)
		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Kind:      d.kind,
			Parameter: d.parameter,
			Payload:   d.payload,
			Evidence:  "Find script: " + evidence,
			Request:   string(d.originalReq),
			Response:  string(dumpedResp),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
		entity.WholeIssue = append(entity.WholeIssue, newIssue)
		entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, err := url.Parse(location)
		if logger.ErrHandle(err) {
			return
		}
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				var err error
				redirectReq, err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
				if logger.ErrHandle(err) {
					return
				}
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq,err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
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
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	var dumpedResp []byte
	b := false
	for _, v := range *d.candidate {
		var inspectReq *http.Request
		var err error
		if len(v.PostParams) != 0 {
			inspectReq, err = genPostParamReq(&v, &v.PostParams)
		} else {
			inspectReq, err = genGetParamReq(&v, &v.GetParams)
		}

		if logger.ErrHandle(err) {
			return
		}

		inspectResp, err := client.Do(inspectReq)
		if logger.ErrHandle(err) {
			continue
		}

		//ここでdumpを行わないとResponseBodyが取れない。
		dumpedResp, err = httputil.DumpResponse(inspectResp, true)
		logger.ErrHandle(err)

		doc, err := goquery.NewDocumentFromReader(inspectResp.Body)
		logger.ErrHandle(err)

		var flg bool
		var evidence string
		doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			injectedPayload := s.Text()
			if strings.Contains(injectedPayload, "alert(\""+d.landmark+"\")") {
				evidence = "alert(\"" + d.landmark + "\")"
				flg = true
				return false
			}
			return true
		})

		if !flg {
			doc.Find("*").EachWithBreak(func(_ int, s *goquery.Selection) bool {
				href, _ := s.Attr("href")
				if strings.HasPrefix(href, "javascript:alert(\""+d.landmark+"\")") {
					evidence = "javascript:alert(\"" + d.landmark + "\")"
					flg = true
					return false
				}
				src, _ := s.Attr("src")
				if strings.HasPrefix(src, "javascript:alert(\""+d.landmark+"\")") {
					evidence = "javascript:alert(\"" + d.landmark + "\")"
					flg = true
					return false
				}
				if src == "x" {
					onerror, _ := s.Attr("onerror")
					if strings.Contains(onerror, "alert(\""+d.landmark+"\")") {
						evidence = "alert(\"" + d.landmark + "\")"
						flg = true
						return false
					}
				}
				onmouseover, _ := s.Attr("onmouseover")
				if strings.Contains(onmouseover, "alert(\""+d.landmark+"\")") {
					evidence = "alert(\"" + d.landmark + "\")"
					flg = true
					return false
				}
				return true
			})
		}

		if flg {
			fmt.Println(d.kind)
			newIssue := entity.Issue{
				URL:       d.jsonMessage.URL,
				Kind:      d.kind,
				Parameter: d.parameter,
				Payload:   d.payload,
				Evidence:  "Find stored script: " + evidence,
				Request:   string(d.originalReq),
				Response:  string(dumpedResp),
			}
			*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
			entity.WholeIssue = append(entity.WholeIssue, newIssue)
			entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
			b = true
		}

		io.ReadAll(inspectResp.Body)
		inspectResp.Body.Close()

		if b {
			return
		}
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()
}

func searchLandmark(d determinant, req []*http.Request) {
	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	var targetResp string
	body, err := io.ReadAll(resp.Body)
	if !logger.ErrHandle(err) {
		targetResp = string(body)
	}

	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, err := url.Parse(location)
		if logger.ErrHandle(err) {
			return
		}
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				var err error
				redirectReq, err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
				if logger.ErrHandle(err) {
					return
				}
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq,err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
		} else {
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		searchLandmark(d, req)
	}

	if strings.Contains(targetResp, d.landmark) {
		// reflect
		return
	} else {
		// stored
		if !QuickScan {
			d.patrol(entity.JsonNodes, d.landmark)
		}
	}
}

func detectHTTPHeaderi(d determinant, req []*http.Request) {
	req[len(req)-1].URL.RawQuery = strings.Replace(req[len(req)-1].URL.RawQuery, "%25", "%", -1)

	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	cookie := resp.Header.Get("Set-Cookie")

	if strings.Contains(cookie, "Himawari=pwned;") {
		dumpedResp, err := httputil.DumpResponse(resp, true)
		logger.ErrHandle(err)

		fmt.Println(d.kind)
		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Kind:      d.kind,
			Parameter: d.parameter,
			Payload:   d.payload,
			Evidence:  "Response Header: " + cookie,
			Request:   string(d.originalReq),
			Response:  string(dumpedResp),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
		entity.WholeIssue = append(entity.WholeIssue, newIssue)
		entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		var redirectReq *http.Request
		l, err := url.Parse(location)
		if logger.ErrHandle(err) {
			return
		}
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				var err error
				redirectReq, err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
				if logger.ErrHandle(err) {
					return
				}
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq,err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
		} else {
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		detectHTTPHeaderi(d, req)
	}
}

func detectCSRF(d determinant, req []*http.Request) {
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	// status code 400, 500番台を排除。もう少し厳しい判定基準や検査対象を絞る必要がある。
	if resp.StatusCode < 400 {
		dumpedResp, err := httputil.DumpResponse(resp, true)
		logger.ErrHandle(err)

		fmt.Println(d.kind)
		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Kind:      d.kind,
			Parameter: d.parameter,
			Payload:   d.payload,
			Evidence:  "Status code: " + resp.Status,
			Request:   string(d.originalReq),
			Response:  string(dumpedResp),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
		entity.WholeIssue = append(entity.WholeIssue, newIssue)
		entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()
}

func detectOpenRedirect(d determinant, req []*http.Request) {
	if loginMsg.URL != "" {
		client.Jar = login(client.Jar)
	}

	var jar4tmp *cookiejar.Jar
	if d.cookie.Name != "" {
		jar4tmp = jar
		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(req[len(req)-1].URL, d.extractCookie(jar4tmp.Cookies(req[len(req)-1].URL)))
	}

	resp, err := client.Do(req[len(req)-1])
	if logger.ErrHandle(err) {
		return
	}

	if jar4tmp != nil {
		client.Jar = jar4tmp
	}

	if len(req) == 1 {
		d.originalReq = logger.DumpedReq
	}

	location := resp.Header.Get("Location")

	//locationがParseできないと検出できないと思うためreturn
	l, err := url.Parse(location)
	if logger.ErrHandle(err) {
		return
	}

	if l.Host == "example.com" {
		dumpedResp, err := httputil.DumpResponse(resp, true)
		if logger.ErrHandle(err) {
			return
		}

		fmt.Println(d.kind)
		newIssue := entity.Issue{
			URL:       d.jsonMessage.URL,
			Kind:      d.kind,
			Parameter: d.parameter,
			Payload:   d.payload,
			Evidence:  "Response header: " + location,
			Request:   string(d.originalReq),
			Response:  string(dumpedResp),
		}
		*d.eachVulnIssue = append(*d.eachVulnIssue, newIssue)
		entity.WholeIssue = append(entity.WholeIssue, newIssue)
		entity.Vulnmap[d.kind].Issues = append(entity.Vulnmap[d.kind].Issues, newIssue)
	}

	io.ReadAll(resp.Body)
	resp.Body.Close()

	if location != "" {
		var redirectReq *http.Request
		l, err := url.Parse(location)
		if logger.ErrHandle(err) {
			return
		}
		redirect := req[len(req)-1].URL.ResolveReference(l)

		if isSameOrigin(req[len(req)-1].URL, redirect) {
			if resp.StatusCode == 301 || resp.StatusCode == 302 {
				var err error
				redirectReq, err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
				if logger.ErrHandle(err) {
					return
				}
			} else {
				return
			}
			/*307リダイレクト時のコード
			if resp.StatusCode == 307 && len(req[len(req)-1].PostForm) != 0 {
				redirectReq = createPostReq(redirect.String(), req[len(req)-1].URL.String(), req[len(req)-1].PostForm)
				redirectReq.PostForm = req[len(req)-1].PostForm
			} else {
				redirectReq,err = createGetReq(redirect.String(), req[len(req)-1].URL.String())
			}
			*/
		} else {
			entity.AppendOutOfOrigin(req[len(req)-1].URL.String(), redirect.String())
			return
		}
		req = append(req, redirectReq)
		detectOpenRedirect(d, req)
	}
}
