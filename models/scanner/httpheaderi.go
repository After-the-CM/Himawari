package scanner

import (
	"Himawari/models/entity"
	"bufio"
)

func HTTPHeaderi(j *entity.JsonNode) {
	d := determinant{
		kind:          HttpHeaderi,
		approach:      detectHTTPHeaderi,
		eachVulnIssue: &j.Issue,
	}

	var payload []string
	p := readfile("models/scanner/payload/" + d.kind + ".txt")
	hhiPayload := bufio.NewScanner(p)
	for hhiPayload.Scan() {
		payload = append(payload, hhiPayload.Text())
	}

	if j.Path == "/" {
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
