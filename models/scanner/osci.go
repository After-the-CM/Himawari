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

	s := SendStruct{
		kind:          OSCI,
		approach:      timeBasedAttack, //ここでapproachを変えられる
		eachVulnIssue: &j.Issue,
	}

	if j.Path == "/" {

		//before
		/*
			if len(j.Messages) == 0 {
				j.Children[0].Messages[0].URL = j.Children[0].URL
				s.jsonMessage = &j.Children[0].Messages[0]
		*/
		if len(j.Messages) != 0 {
			u, _ := url.Parse(j.URL)
			slash, _ := url.Parse("/")
			j.Messages[0].URL = u.ResolveReference(slash).String()
			s.jsonMessage = &j.Messages[0]

		} else {
			for i, v := range j.Children {
				if len(v.Messages) != 0 {
					j.Children[i].Messages[0].URL = j.Children[i].URL
					s.jsonMessage = &j.Children[i].Messages[0]
					continue
				}
			}
		}
		for _, v := range payload {
			s.setHeaderDocumentRoot(v)
		}

	}

	for i := 0; i < len(j.Messages); i++ {
		j.Messages[i].URL = j.URL
		for _, v := range payload {
			s.jsonMessage = &j.Messages[i]
			s.setParam(v)
			if len(j.Messages[i].PostParams) != 0 {
				s.setPostHeader(v)
			} else {
				s.setGetHeader(v)
			}
		}
	}
	entity.WholeIssue = append(entity.WholeIssue, *s.eachVulnIssue...)
}
