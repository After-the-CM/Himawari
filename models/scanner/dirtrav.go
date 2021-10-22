package scanner

import (
	"Himawari/models/entity"
	"bufio"
)

func DirTrav(j *entity.JsonNode) {
	payload := make([]string, 0, 20)
	p := readfile("models/scanner/payload/dirtrav.txt")
	dirtravPayload := bufio.NewScanner(p)
	for dirtravPayload.Scan() {
		payload = append(payload, dirtravPayload.Text())
	}

	d := determinant{
		kind:          dirTrav,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
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
