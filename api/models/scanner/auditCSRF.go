package scanner

import (
	"fmt"

	"Himawari/models/entity"
)

func auditCSRF(j *entity.JsonNode) {
	fmt.Printf("\x1b[36m%s\x1b[0m\n", "CSRFの診断を開始しました")
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
