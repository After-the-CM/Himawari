package scanner

import (
	"Himawari/models/entity"
)

func Scan(j *entity.JsonNode) {
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
	entity.WholeIssue = []entity.Issue{}
}
