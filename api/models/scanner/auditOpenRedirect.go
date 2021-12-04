package scanner

import (
	"bufio"
	"fmt"

	"Himawari/models/entity"
)

func auditOpenRedirect(j *entity.JsonNode) {

	fmt.Println("\x1b[32m"+"Scan OpenRedirect", "\x1b[0m")

	d := determinant{
		kind:          openRedirect,
		approach:      detectOpenRedirect,
		eachVulnIssue: &j.Issue,
	}

	var payload []string
	p := readfile("models/scanner/payload/" + d.kind + ".txt")
	openRedirectPayload := bufio.NewScanner(p)
	for openRedirectPayload.Scan() {
		payload = append(payload, openRedirectPayload.Text())
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
