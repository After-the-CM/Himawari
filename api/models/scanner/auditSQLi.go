package scanner

import (
	"bufio"
	"fmt"

	"Himawari/models/entity"
)

func auditSQLi(j *entity.JsonNode) {
	e := determinant{
		kind:          errBasedSQLi,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	t := determinant{
		kind:          timeBasedSQLi,
		approach:      timeBasedAttack,
		eachVulnIssue: &j.Issue,
	}

	fmt.Printf("\x1b[36m%s%s%s\x1b[0m\n", "ð", e.kind, "ãŪčĻšæ­ãéå§ããūããð")

	var errSQLiPayloads []string
	ep := readfile("models/scanner/payload/" + e.kind + ".txt")
	ePayloads := bufio.NewScanner(ep)
	for ePayloads.Scan() {
		errSQLiPayloads = append(errSQLiPayloads, ePayloads.Text())
	}

	var timeSQLiPayloads []string
	tp := readfile("models/scanner/payload/" + t.kind + ".txt")
	tPayloads := bufio.NewScanner(tp)
	for tPayloads.Scan() {
		timeSQLiPayloads = append(timeSQLiPayloads, tPayloads.Text())
	}

	var vulnNum int
	if j.Path == "/" {
		vulnNum = len(*e.eachVulnIssue)
		for _, v := range j.Children {
			e.jsonMessage = retrieveJsonMessage(&v)
			t.jsonMessage = e.jsonMessage
			if e.jsonMessage != nil {
				for _, v := range errSQLiPayloads {
					e.setHeaderDocumentRoot(v)
				}
				if vulnNum-len(*e.eachVulnIssue) == 0 {
					for _, v := range timeSQLiPayloads {
						t.setHeaderDocumentRoot(v)
					}
				}
				break
			}
		}
	}

	for i := 0; i < len(j.Messages); i++ {
		vulnNum = len(*e.eachVulnIssue)
		for _, v := range errSQLiPayloads {
			e.jsonMessage = &j.Messages[i]
			e.setParam(v)
			for _, cookie := range j.Cookies {
				e.setCookie(cookie, v)
			}
			if len(j.Messages[i].PostParams) != 0 {
				e.setPostHeader(v)
			} else {
				e.setGetHeader(v)
			}
		}

		if len(*e.eachVulnIssue)-vulnNum != 0 {
			return
		}

		for _, v := range timeSQLiPayloads {
			t.jsonMessage = &j.Messages[i]
			t.setParam(v)
			for _, cookie := range j.Cookies {
				t.setCookie(cookie, v)
			}
			if len(j.Messages[i].PostParams) != 0 {
				t.setPostHeader(v)
			} else {
				t.setGetHeader(v)
			}
		}

	}
}
