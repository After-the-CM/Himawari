package scanner

import (
	"fmt"
	"net/http"

	"Himawari/models/entity"
	"Himawari/models/logger"
)

func auditDirListing(j *entity.JsonNode) {
	d := determinant{
		kind:          dirListing,
		approach:      stringMatching,
		eachVulnIssue: &j.Issue,
	}

	fmt.Printf("\x1b[36m%s%s%s\x1b[0m\n", "🔍", d.kind, "の診断を開始しました🔍")

	req, err := createGetReq(j.URL, "")
	if logger.ErrHandle(err) {
		return
	}

	//最初は`/`がついていないURL
	d.approach(d, []*http.Request{req})

	reqslash, err := createGetReq(j.URL+"/", "")
	if logger.ErrHandle(err) {
		return
	}
	d.approach(d, []*http.Request{reqslash})

}
