package scanner

import (
	"bufio"
	"fmt"

	"Himawari/models/entity"
)

func auditDirTraversal(j *entity.JsonNode) {
	d := determinant{
		kind:          dirTraversal,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	fmt.Printf("\x1b[36m%s%s%s\x1b[0m\n", "ð", d.kind, "ãŪčĻšæ­ãéå§ããūããð")

	var payload []string
	p := readfile("models/scanner/payload/" + d.kind + ".txt")
	dirtravPayload := bufio.NewScanner(p)
	for dirtravPayload.Scan() {
		payload = append(payload, dirtravPayload.Text())
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
			for _, cookie := range j.Cookies {
				d.setCookie(cookie, v)
			}
			if len(j.Messages[i].PostParams) != 0 {
				d.setPostHeader(v)
			} else {
				d.setGetHeader(v)
			}
		}

	}
}
