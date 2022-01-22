package scanner

import (
	"fmt"

	"Himawari/models/entity"
)

func Scan(j *entity.JsonNode) {
	fmt.Printf("\x1b[36m%s%s%s\x1b[0m\n", "🌻", j.URL, "への診断を開始します🌻")
	auditOSCi(j)
	auditDirTraversal(j)
	auditSQLi(j)
	auditOpenRedirect(j)
	auditDirListing(j)
	auditXSS(j)
	auditHTTPHeaderi(j)
	auditCSRF(j)

	if len(j.Children) > 0 {
		for i := 0; i < len(j.Children); i++ {
			Scan(&j.Children[i])
		}
	}
}

func Reset() {
	entity.ResetVulnMap()
	entity.WholeIssue = []entity.Issue{}
	entity.OutOfOrigin = map[string][]string{}
}
