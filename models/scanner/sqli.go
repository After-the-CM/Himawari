package scanner

import (
	"bufio"

	"Himawari/models/entity"
)

func SQLi(j *entity.JsonNode) {
	e := determinant{
		kind:          ErrBasedSQLi,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	t := determinant{
		kind:          TimeBasedSQLi,
		approach:      timeBasedAttack,
		eachVulnIssue: &j.Issue,
	}

	var errSQLiPayloads []string
	ePayload := readfile("models/scanner/payload/" + e.kind + ".txt")
	ePayloads := bufio.NewScanner(ePayload)
	for ePayloads.Scan() {
		errSQLiPayloads = append(errSQLiPayloads, ePayloads.Text())
	}

	var timeSQLiPayloads []string
	tPayload := readfile("models/scanner/payload/" + t.kind + ".txt")
	tPayloads := bufio.NewScanner(tPayload)
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
			if len(j.Messages[i].PostParams) != 0 {
				t.setPostHeader(v)
			} else {
				t.setGetHeader(v)
			}
		}

	}
}
