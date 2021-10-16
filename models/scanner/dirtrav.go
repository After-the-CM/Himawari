package scanner

import (
	"Himawari/models/entity"
	"bufio"
	"net/url"
)

func DirTrav(j *entity.JsonNode) {
	payload := make([]string, 0, 20)
	p := readfile("models/scanner/payload/dirtrav.txt")
	osciPayload := bufio.NewScanner(p)
	for osciPayload.Scan() {
		payload = append(payload, osciPayload.Text())
	}

	d := determinant{
		kind:          dirTrav,
		approach:      stringMatching, //ここでsendMethodを変えられる
		eachVulnIssue: &j.Issue,
	}

	if j.Path == "/" {
		if len(j.Messages) != 0 {
			//crawl時に入力されたURLの語尾に`/`がない場合の対処
			u, _ := url.Parse(j.Messages[0].URL)
			slash, _ := url.Parse("/")
			j.Messages[0].URL = u.ResolveReference(slash).String()
			d.jsonMessage = &j.Messages[0]
		} else {
			for i, v := range j.Children {
				if len(v.Messages) != 0 {
					d.jsonMessage = &j.Children[i].Messages[0]
					continue
				}
			}
		}
		for _, v := range payload {
			d.setHeaderDocumentRoot(v)
		}

	}

	//	for i, v := range j.Messages {
	for i := 0; i < len(j.Messages); i++ {
		//j.Messages[i].URL = j.URL
		for _, v := range payload {
			d.jsonMessage = &j.Messages[i]
			d.setParam(v)
			if len(j.Messages[i].PostParams) != 0 {
				d.setPostHeader(v)
			} else {
				d.setGetHeader(v)
			}
		}

	}
	//appendできているかな？無理そうならポインタに
	//j.Issue = append(j.Issue, nodeIssue...)

	//j.Issue = append(j.Issue, *s.eachVulnIssue...)
	//Issues = append(Issues, nodeIssue...)

	//entity.WholeIssue = append(entity.WholeIssue, *s.eachVulnIssue...)
	//return j.Issue
}
