package scanner

import (
	"fmt"

	"Himawari/models/entity"
)

func auditCSRF(j *entity.JsonNode) {
	d := determinant{
		kind:          csrf,
		approach:      detectCSRF,
		eachVulnIssue: &j.Issue,
	}

	fmt.Printf("\x1b[36m%s%s%s\x1b[0m\n", "ð", d.kind, "ãŪčĻšæ­ãéå§ããūããð")

	for i := 0; i < len(j.Messages); i++ {
		d.jsonMessage = &j.Messages[i]
		if len(j.Messages[i].PostParams) != 0 {
			d.setPostRef("")
		}
	}
}
