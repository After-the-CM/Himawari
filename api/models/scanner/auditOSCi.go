package scanner

import (
	"bufio"
	"fmt"

	"Himawari/models/entity"
)

func auditOSCi(j *entity.JsonNode) {

	fmt.Println("\x1b[32m"+"Scan OSCI\n", "\x1b[0m")

	d := determinant{
		kind:          osci,
		approach:      timeBasedAttack, //ここでapproachを変えられる
		eachVulnIssue: &j.Issue,
	}

	var payload []string
	p := readfile("models/scanner/payload/" + d.kind + ".txt")
	osciPayload := bufio.NewScanner(p)
	for osciPayload.Scan() {
		payload = append(payload, osciPayload.Text())
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
