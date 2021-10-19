package scanner

import (
	"bufio"
	"net/url"

	"Himawari/models/entity"
)

func Osci(j *entity.JsonNode) {
	payload := make([]string, 0, 20)
	p := readfile("models/scanner/payload/osci.txt")
	osciPayload := bufio.NewScanner(p)
	for osciPayload.Scan() {
		payload = append(payload, osciPayload.Text())
	}

	d := determinant{
		kind:          OSCI,
		approach:      timeBasedAttack, //ここでapproachを変えられる
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

	for i := 0; i < len(j.Messages); i++ {
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
}
