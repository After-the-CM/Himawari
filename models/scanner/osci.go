package scanner

import (
	"bufio"

	"Himawari/models/entity"
)

func Osci(j *entity.JsonNode) {
	payload := make([]string, 0, 20)
	p := readfile("models/scanner/payload/" + OSCI + ".txt")
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
		//直接retrieveJsonMessageにjを渡してもよいが、`/`問題を解決するために、forで回している。
		for _, v := range j.Children {
			d.jsonMessage = retrieveJsonMessage(&v)
			if d.jsonMessage != nil {
				for _, v := range payload {
					d.setHeaderDocumentRoot(v)
				}
				break
			}
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
