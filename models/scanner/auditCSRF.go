package scanner

import (
	"Himawari/models/entity"
)

func auditCSRF(j *entity.JsonNode) {

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
