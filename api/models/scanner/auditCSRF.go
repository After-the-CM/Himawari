package scanner

import (
	"Himawari/models/entity"
	"fmt"
)

func auditCSRF(j *entity.JsonNode) {

	fmt.Println("\x1b[32m"+"Scan CSRF", "\x1b[0m")

	d := determinant{
		kind:          csrf,
		approach:      detectCSRF,
		eachVulnIssue: &j.Issue,
	}

	for i := 0; i < len(j.Messages); i++ {
		d.jsonMessage = &j.Messages[i]
		if len(j.Messages[i].PostParams) != 0 {
			d.setPostRef("")
		}
	}
}
