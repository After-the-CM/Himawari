package scanner

import (
	"Himawari/models/entity"
	"bufio"
)

func SQLi(j *entity.JsonNode) {
	errSQLiPayloads := make([]string, 0, 1)
	ePayload := readfile("models/scanner/payload/errbasedsqli.txt")
	ePayloads := bufio.NewScanner(ePayload)
	for ePayloads.Scan() {
		errSQLiPayloads = append(errSQLiPayloads, ePayloads.Text())
	}

	timeSQLiPayloads := make([]string, 0, 3)
	tPayload := readfile("models/scanner/payload/timebasedsqli.txt")
	tPayloads := bufio.NewScanner(tPayload)
	for tPayloads.Scan() {
		timeSQLiPayloads = append(timeSQLiPayloads, tPayloads.Text())
	}

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

	var vulnNum int
	if j.Path == "/" {
		vulnNum = len(*e.eachVulnIssue)
		for _, v := range j.Children {
			e.jsonMessage = retrieveJsonMessage(&v)
			t.jsonMessage = retrieveJsonMessage(&v)
			if e.jsonMessage != nil {
				for _, v := range errSQLiPayloads {
					e.setHeaderDocumentRoot(v)
				}
				if vulnNum-len(*e.eachVulnIssue) != 0 {
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
