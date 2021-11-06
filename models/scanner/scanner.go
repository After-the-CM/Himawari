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

//再度scanを実行する際に各jsonNodeのIssueが埋まっている状態だと動かないため。
/*
func ResetJsonMessageOfIssue() {
	entity.JsonNodes.Issue = []entity.Issue{}
	resetJsonMessageOfIssue(&entity.JsonNodes)
}

func resetJsonMessageOfIssue(j *entity.JsonNode) {
	for i, v := range j.Children {
		if len(v.Children) != 0 {
			resetJsonMessageOfIssue(&j.Children[i])
		}
		j.Children[i].Issue = []entity.Issue{}
	}
}
*/
