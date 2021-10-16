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
		//kind: SQLI,
		kind:          ErrBasedSQLi,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	t := determinant{
		//kind: SQLI,
		kind:          TimeBasedSQLi,
		approach:      timeBasedAttack,
		eachVulnIssue: &j.Issue,
	}

	var vulnNum int

	if j.Path == "/" {
		vulnNum = len(*e.eachVulnIssue)
		if len(j.Messages) == 0 {
			j.Children[0].Messages[0].URL = j.Children[0].URL
			e.jsonMessage = &j.Children[0].Messages[0]
			t.jsonMessage = &j.Children[0].Messages[0]
		} else {
			for i, v := range j.Children {
				if len(v.Messages) != 0 {
					j.Children[i].Messages[0].URL = j.Children[i].URL
					e.jsonMessage = &j.Children[i].Messages[0]
					t.jsonMessage = &j.Children[i].Messages[0]
					continue
				}
			}
			for _, v := range errSQLiPayloads {
				e.setHeaderDocumentRoot(v)
			}
			// error basedのpayloadが刺さらなかったら
			if vulnNum-len(*e.eachVulnIssue) == 0 {
				for _, v := range timeSQLiPayloads {
					t.setHeaderDocumentRoot(v)
				}
			}
		}
	}
	for i := 0; i < len(j.Messages); i++ {
		j.Messages[i].URL = j.URL
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
		// error basedのpayloadが刺さらなかったら
		if len(*e.eachVulnIssue)-vulnNum == 0 {
			entity.WholeIssue = append(entity.WholeIssue, *e.eachVulnIssue...)
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
	entity.WholeIssue = append(entity.WholeIssue, *t.eachVulnIssue...)
}